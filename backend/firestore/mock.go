package firestore

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
)

// CollectionRefInterface defines the interface for collection operations
type CollectionRefInterface interface {
	Documents(ctx context.Context) DocumentIteratorInterface
	Where(field, op string, value interface{}) CollectionRefInterface
	Limit(n int) CollectionRefInterface
	Add(ctx context.Context, data interface{}) (DocumentRefInterface, *firestore.WriteResult, error)
	Doc(id string) DocumentRefInterface
}

// DocumentIteratorInterface defines the interface for document iteration
type DocumentIteratorInterface interface {
	GetAll() ([]DocumentSnapshotInterface, error)
}

// DocumentSnapshotInterface defines the interface for document snapshots
type DocumentSnapshotInterface interface {
	DataTo(dest interface{}) error
	GetID() string
	GetRef() DocumentRefInterface
}

// DocumentRefInterface defines the interface for document references
type DocumentRefInterface interface {
	Set(ctx context.Context, data interface{}) (*firestore.WriteResult, error)
	Collection(name string) CollectionRefInterface
	GetID() string
}

// RealCollectionRef wraps a real Firestore collection reference
type RealCollectionRef struct {
	ref   *firestore.CollectionRef
	query firestore.Query
	hasQuery bool
}

func (r *RealCollectionRef) Documents(ctx context.Context) DocumentIteratorInterface {
	if r.hasQuery {
		return &RealDocumentIterator{iter: r.query.Documents(ctx)}
	}
	return &RealDocumentIterator{iter: r.ref.Documents(ctx)}
}

func (r *RealCollectionRef) Where(field, op string, value interface{}) CollectionRefInterface {
	var q firestore.Query
	if r.hasQuery {
		q = r.query.Where(field, op, value)
	} else {
		q = r.ref.Where(field, op, value)
	}
	return &RealCollectionRef{ref: r.ref, query: q, hasQuery: true}
}

func (r *RealCollectionRef) Limit(n int) CollectionRefInterface {
	var q firestore.Query
	if r.hasQuery {
		q = r.query.Limit(n)
	} else {
		q = r.ref.Limit(n)
	}
	return &RealCollectionRef{ref: r.ref, query: q, hasQuery: true}
}

func (r *RealCollectionRef) Add(ctx context.Context, data interface{}) (DocumentRefInterface, *firestore.WriteResult, error) {
	docRef, wr, err := r.ref.Add(ctx, data)
	if err != nil {
		return nil, nil, err
	}
	return &RealDocumentRef{ref: docRef}, wr, nil
}

func (r *RealCollectionRef) Doc(id string) DocumentRefInterface {
	return &RealDocumentRef{ref: r.ref.Doc(id)}
}

// RealDocumentIterator wraps a real Firestore document iterator
type RealDocumentIterator struct {
	iter *firestore.DocumentIterator
}

func (r *RealDocumentIterator) GetAll() ([]DocumentSnapshotInterface, error) {
	docs, err := r.iter.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]DocumentSnapshotInterface, len(docs))
	for i, doc := range docs {
		result[i] = &RealDocumentSnapshot{snapshot: doc}
	}
	return result, nil
}

// RealDocumentSnapshot wraps a real Firestore document snapshot
type RealDocumentSnapshot struct {
	snapshot *firestore.DocumentSnapshot
}

func (r *RealDocumentSnapshot) DataTo(dest interface{}) error {
	return r.snapshot.DataTo(dest)
}

func (r *RealDocumentSnapshot) GetID() string {
	return r.snapshot.Ref.ID
}

func (r *RealDocumentSnapshot) GetRef() DocumentRefInterface {
	return &RealDocumentRef{ref: r.snapshot.Ref}
}

// RealDocumentRef wraps a real Firestore document reference
type RealDocumentRef struct {
	ref *firestore.DocumentRef
}

func (r *RealDocumentRef) Set(ctx context.Context, data interface{}) (*firestore.WriteResult, error) {
	return r.ref.Set(ctx, data)
}

func (r *RealDocumentRef) Collection(name string) CollectionRefInterface {
	return &RealCollectionRef{ref: r.ref.Collection(name)}
}

func (r *RealDocumentRef) GetID() string {
	return r.ref.ID
}

// MockFirestoreClient is an in-memory mock implementation of Firestore
type MockFirestoreClient struct {
	mu         sync.RWMutex
	attendees  map[string]map[string]interface{}
	sessions   map[string]map[string]interface{}
	speakers   map[string]map[string]interface{}
	nextID     int
	clientID   string
}

// MockDocumentSnapshot represents a Firestore document snapshot
type MockDocumentSnapshot struct {
	ID   string
	Data map[string]interface{}
	Ref  *MockDocumentRef
}

// MockDocumentRef represents a Firestore document reference
type MockDocumentRef struct {
	ID       string
	Path     string
	Parent   *MockCollectionRef
	Client   *MockFirestoreClient
}

// MockCollectionRef represents a Firestore collection reference
type MockCollectionRef struct {
	ID     string
	Path   string
	Parent *MockDocumentRef
	Client *MockFirestoreClient
	Data   *map[string]map[string]interface{}
}

// MockQueryIterator represents a Firestore query iterator
type MockQueryIterator struct {
	Docs []*MockDocumentSnapshot
	Idx  int
}

// Implement CollectionRefInterface for MockCollectionRef
func (c *MockCollectionRef) Documents(ctx context.Context) DocumentIteratorInterface {
	return &MockQueryIterator{Docs: c.getAllDocs()}
}

func (c *MockCollectionRef) Where(field, op string, value interface{}) CollectionRefInterface {
	filtered := c.filterDocs(field, op, value)
	return &MockCollectionRef{
		ID:     c.ID,
		Path:   c.Path,
		Parent: c.Parent,
		Client: c.Client,
		Data:   &filtered,
	}
}

func (c *MockCollectionRef) Limit(n int) CollectionRefInterface {
	limited := c.limitDocs(n)
	return &MockCollectionRef{
		ID:     c.ID,
		Path:   c.Path,
		Parent: c.Parent,
		Client: c.Client,
		Data:   &limited,
	}
}

func (c *MockCollectionRef) Add(ctx context.Context, data interface{}) (DocumentRefInterface, *firestore.WriteResult, error) {
	if c.Data == nil {
		return nil, nil, fmt.Errorf("collection not initialized")
	}
	
	c.Client.mu.Lock()
	defer c.Client.mu.Unlock()
	
	id := fmt.Sprintf("doc_%d_%d", c.Client.nextID, time.Now().UnixNano())
	c.Client.nextID++
	
	dataMap := structToMap(data)
	(*c.Data)[id] = dataMap
	
	docRef := &MockDocumentRef{
		ID:     id,
		Path:   c.Path + "/" + id,
		Parent: c,
		Client: c.Client,
	}
	
	return docRef, &firestore.WriteResult{}, nil
}

func (c *MockCollectionRef) Doc(id string) DocumentRefInterface {
	return &MockDocumentRef{
		ID:     id,
		Path:   c.Path + "/" + id,
		Parent: c,
		Client: c.Client,
	}
}

// Implement DocumentIteratorInterface for MockQueryIterator
func (it *MockQueryIterator) GetAll() ([]DocumentSnapshotInterface, error) {
	result := make([]DocumentSnapshotInterface, len(it.Docs))
	for i, doc := range it.Docs {
		result[i] = doc
	}
	return result, nil
}

// Implement DocumentSnapshotInterface for MockDocumentSnapshot
func (d *MockDocumentSnapshot) DataTo(dest interface{}) error {
	return mapToStruct(d.Data, dest)
}

func (d *MockDocumentSnapshot) GetID() string {
	return d.ID
}

func (d *MockDocumentSnapshot) GetRef() DocumentRefInterface {
	return d.Ref
}

// Implement DocumentRefInterface for MockDocumentRef
func (d *MockDocumentRef) Set(ctx context.Context, data interface{}) (*firestore.WriteResult, error) {
	if d.Parent == nil || d.Parent.Data == nil {
		return nil, fmt.Errorf("parent collection not initialized")
	}
	
	d.Client.mu.Lock()
	defer d.Client.mu.Unlock()
	
	dataMap := structToMap(data)
	(*d.Parent.Data)[d.ID] = dataMap
	
	return &firestore.WriteResult{}, nil
}

func (d *MockDocumentRef) Collection(name string) CollectionRefInterface {
	var dataMap *map[string]map[string]interface{}
	switch name {
	case "attendees":
		dataMap = &d.Client.attendees
	case "sessions":
		dataMap = &d.Client.sessions
	case "speakers":
		dataMap = &d.Client.speakers
	default:
		dataMap = &map[string]map[string]interface{}{}
	}
	
	return &MockCollectionRef{
		ID:     name,
		Path:   d.Path + "/" + name,
		Parent: d,
		Client: d.Client,
		Data:   dataMap,
	}
}

func (d *MockDocumentRef) GetID() string {
	return d.ID
}

// Helper methods for MockCollectionRef
func (c *MockCollectionRef) getAllDocs() []*MockDocumentSnapshot {
	if c.Data == nil {
		return []*MockDocumentSnapshot{}
	}
	
	c.Client.mu.RLock()
	defer c.Client.mu.RUnlock()
	
	docs := make([]*MockDocumentSnapshot, 0)
	for id, data := range *c.Data {
		docs = append(docs, &MockDocumentSnapshot{
			ID:   id,
			Data: data,
			Ref: &MockDocumentRef{
				ID:     id,
				Path:   c.Path + "/" + id,
				Parent: c,
				Client: c.Client,
			},
		})
	}
	return docs
}

func (c *MockCollectionRef) filterDocs(field, op string, value interface{}) map[string]map[string]interface{} {
	if c.Data == nil {
		return make(map[string]map[string]interface{})
	}
	
	c.Client.mu.RLock()
	defer c.Client.mu.RUnlock()
	
	filtered := make(map[string]map[string]interface{})
	for id, data := range *c.Data {
		if val, ok := data[field]; ok {
			if op == "==" && reflect.DeepEqual(val, value) {
				filtered[id] = data
			}
		}
	}
	return filtered
}

func (c *MockCollectionRef) limitDocs(n int) map[string]map[string]interface{} {
	if c.Data == nil {
		return make(map[string]map[string]interface{})
	}
	
	c.Client.mu.RLock()
	defer c.Client.mu.RUnlock()
	
	limited := make(map[string]map[string]interface{})
	count := 0
	for id, data := range *c.Data {
		if count >= n {
			break
		}
		limited[id] = data
		count++
	}
	return limited
}

// Collection returns a collection reference
func (m *MockFirestoreClient) Collection(name string) CollectionRefInterface {
	if name == "clients" {
		return &MockCollectionRef{
			ID:     name,
			Path:   name,
			Client: m,
		}
	}
	return nil
}

// Doc returns a document reference for the clients collection
func (m *MockFirestoreClient) Doc(id string) DocumentRefInterface {
	return &MockDocumentRef{
		ID:     id,
		Path:   "clients/" + id,
		Client: m,
		Parent: &MockCollectionRef{
			ID:     "clients",
			Path:   "clients",
			Client: m,
		},
	}
}

// NewMockFirestoreClient creates a new mock Firestore client
func NewMockFirestoreClient(clientID string) *MockFirestoreClient {
	return &MockFirestoreClient{
		attendees: make(map[string]map[string]interface{}),
		sessions:  make(map[string]map[string]interface{}),
		speakers:  make(map[string]map[string]interface{}),
		nextID:    1,
		clientID:  clientID,
	}
}

// Helper function to convert struct to map
func structToMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return result
	}
	
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		
		tag := field.Tag.Get("firestore")
		if tag == "" || tag == "-" {
			tag = field.Tag.Get("json")
			if tag == "" || tag == "-" {
				tag = field.Name
			}
		}
		
		if idx := len(tag); idx > 0 {
			for i, r := range tag {
				if r == ',' {
					idx = i
					break
				}
			}
			tag = tag[:idx]
		}
		
		if fieldVal.CanInterface() {
			result[tag] = fieldVal.Interface()
		}
	}
	
	return result
}

// Helper function to convert map to struct
func mapToStruct(m map[string]interface{}, dest interface{}) error {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}
	
	val = val.Elem()
	typ := val.Type()
	
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to struct")
	}
	
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		
		tag := field.Tag.Get("firestore")
		if tag == "" || tag == "-" {
			tag = field.Tag.Get("json")
			if tag == "" || tag == "-" {
				tag = field.Name
			}
		}
		
		if idx := len(tag); idx > 0 {
			for i, r := range tag {
				if r == ',' {
					idx = i
					break
				}
			}
			tag = tag[:idx]
		}
		
		if value, ok := m[tag]; ok && fieldVal.CanSet() {
			setValue(fieldVal, value)
		}
	}
	
	return nil
}

// Helper function to set a value
func setValue(field reflect.Value, value interface{}) {
	if !field.CanSet() {
		return
	}
	
	val := reflect.ValueOf(value)
	
	if field.Type() == reflect.TypeOf(time.Time{}) {
		if t, ok := value.(time.Time); ok {
			field.Set(reflect.ValueOf(t))
			return
		}
		if t, ok := value.(*time.Time); ok && t != nil {
			field.Set(reflect.ValueOf(*t))
			return
		}
		if str, ok := value.(string); ok {
			if t, err := time.Parse(time.RFC3339, str); err == nil {
				field.Set(reflect.ValueOf(t))
				return
			}
		}
	}
	
	if val.Type().AssignableTo(field.Type()) {
		field.Set(val)
		return
	}
	
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
		return
	}
}
