// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

// Executed before an object is initially persisted to the database
type BeforeInsertProcessor interface {
	BeforeInsert()
}

// Executed before an object is updated
type BeforeUpdateProcessor interface {
	BeforeUpdate()
}

// Executed before an object is deleted
type BeforeDeleteProcessor interface {
	BeforeDelete()
}

type BeforeSetProcessor interface {
	BeforeSet(string, Cell)
}

type AfterSetProcessor interface {
	AfterSet(string, Cell)
}

// !nashtsai! TODO enable BeforeValidateProcessor when xorm start to support validations
//// Executed before an object is validated
//type BeforeValidateProcessor interface {
//    BeforeValidate()
//}
// --

// Executed after an object is persisted to the database
type AfterInsertProcessor interface {
	AfterInsert()
}

// Executed after an object has been updated
type AfterUpdateProcessor interface {
	AfterUpdate()
}

// Executed after an object has been deleted
type AfterDeleteProcessor interface {
	AfterDelete()
}
