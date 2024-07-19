package service

// var _ BlobStore = (*testMock)(nil)
// var _ BucketStore = (*testMock)(nil)
// var _ AccountStore = (*testMock)(nil)

// type testMock struct {
// 	accountStore map[string]*user.Account
// 	blobStore    map[string]string
// 	BucketStore  map[string]map[string]*blob.ObjectInfo
// }

// // Delete implements BlobStore.
// func (t *testMock) Delete(ctx context.Context, bucket string, filename string) error {
// b,ok := t.BucketStore[bucket]
// if !ok {
// 	return nil
// }
// 	delete(t.blobStore[b], )
// }

// // Find implements AccountStore.
// func (t *testMock) Find(ctx context.Context, username string) (*user.Account, error) {
// 	// panic("unimplemented")
// 	acc, ok := t.accountStore[username]
// 	if !ok {
// 		return nil, fmt.Errorf("invalid username")
// 	}
// 	return acc, nil
// }

// // Insert implements AccountStore.
// func (t *testMock) Insert(ctx context.Context, acc *user.Account) error {
// 	t.accountStore[acc.Name] = acc
// 	return nil
// }

// // CreateBucket implements BucketStore.
// func (t *testMock) CreateBucket(ctx context.Context, bucket string) (any, error) {
// 	t.BucketStore[bucket] = map[string]*blob.ObjectInfo{}
// 	return nil, nil
// }

// // IsBucketExist implements BucketStore.
// func (t *testMock) IsBucketExist(ctx context.Context, bucket string) (bool, error) {
// 	_, ok := t.BucketStore[bucket]
// 	return ok, nil
// }

// // ObjectIterator implements BucketStore.
// func (t *testMock) ObjectIterator(ctx context.Context, bucket string, limit int, lastFilename string) *blob.Iterator {
// 	blobs := t.BucketStore[bucket]

// 	if len(blobs) <= 0 {
// 		limit = 1
// 	}
// 	objC := make(chan *blob.ObjectInfo, limit)
// 	for _, v := range blobs {
// 		objC <- v
// 	}
// 	close(objC)

// 	iter := blob.Iterator{
// 		UserID: "",
// 		C:      objC,
// 	}

// 	return &iter
// }

// // Get implements BlobStore.
// func (t *testMock) Get(ctx context.Context, userID string, filename string) (*blob.ObjectInfo, error) {
// 	bs := t.BucketStore[userID]
// 	obj := bs[filename]
// 	data := t.blobStore[userID+filename]

// 	obj.Reader = blob.ReaderFunc(func(ctx context.Context) (io.ReadCloser, error) {
// 		rc := io.NopCloser(bytes.NewBuffer([]byte(data)))
// 		return rc, nil
// 	})

// 	return obj, nil
// }

// // Put implements BlobStore.
// func (t *testMock) Put(ctx context.Context, userID string, filename string, reader io.Reader, size int64, contentType string) (*blob.ObjectInfo, error) {
// 	store := t.BucketStore[userID]
// 	obj := &blob.ObjectInfo{
// 		UserID:       userID,
// 		Filename:     filename,
// 		ContentType:  contentType,
// 		Sum:          "",
// 		Size:         size,
// 		LastModified: time.Now(),
// 	}
// 	store[filename] = obj

// 	data, _ := io.ReadAll(reader)
// 	t.blobStore[userID+filename] = string(data)

// 	return obj, nil
// }

// func Test_app_accountService(t *testing.T) {
// 	tMock := testMock{
// 		accountStore: map[string]*user.Account{},
// 		// blobStore:    map[string]*blob.ObjectInfo{},
// 		BucketStore: map[string]map[string]*blob.ObjectInfo{},
// 	}

// 	app, _ := NewCloudfs(&tMock, &tMock, &tMock)
// 	if err := app.Register(context.Background(), &RegisterParam{Username: "uname", Password: "pass123"}); err != nil {
// 		t.Fatal(err)
// 	}
// 	if _, err := app.Auth(context.Background(), &AuthParam{
// 		Username: "uname",
// 		Password: "pass123",
// 	}); err != nil {
// 		t.Fatal(err)
// 	}

// }

// func Test_app_blob(t *testing.T) {
// 	tMock := testMock{
// 		accountStore: map[string]*user.Account{},
// 		blobStore:    map[string]string{},
// 		BucketStore:  map[string]map[string]*blob.ObjectInfo{},
// 	}

// 	UserID := "1234"
// 	ObjName := "my-object"
// 	ContentType := "content"
// 	data := "content of object"
// 	buf := bytes.NewBuffer([]byte(data))

// 	app, _ := NewCloudfs(&tMock, &tMock, &tMock)
// 	// if _, err := app.MakeBucket(context.Background(), expectBucket); err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	if _, err := app.Upload(context.Background(), &UploadParam{
// 		UserID:      UserID,
// 		Filename:    ObjName,
// 		Size:        int64(len(data)),
// 		ContentType: ContentType,
// 		DataReader:  buf,
// 	}); err != nil {
// 		t.Fatal(err)
// 	}

// 	infos, err := app.ListObject(context.Background(), UserID, 1, "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	assert.Len(t, infos, 1)
// 	for _, actual := range infos {
// 		assert.Equal(t, UserID, actual.UserID)
// 	}

// 	_, err = app.Download(context.Background(), buf, &DownloadParam{
// 		UserID:   UserID,
// 		Filename: ObjName,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.Equal(t, string(data), buf.String())
// }

// func Test_err(t *testing.T) {
// 	up := UploadParam{
// 		UserID:      "",
// 		Filename:    "",
// 		Size:        0,
// 		ContentType: "",
// 		DataReader:  nil,
// 	}

// 	err := up.validate()
// 	if !errors.Is(err, ErrUpload) {
// 		t.Fatal()
// 	}
// }
