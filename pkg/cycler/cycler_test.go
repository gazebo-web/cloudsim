package cycler

import (
    "github.com/stretchr/testify/suite"
    "testing"
)

func TestCyclerSuite(t *testing.T) {
    suite.Run(t, &CyclerTestSuite{})
}

type CyclerTestSuite struct {
    suite.Suite
}

func (s *CyclerTestSuite) TestGet() {
    values := []interface{}{1, 2, 3}
    cycler := &cycler{
        values: values,
    }
    // Get should always return the same value
    s.Require().Equal(1, cycler.Get())
    s.Require().Equal(1, cycler.Get())
    s.Require().Equal(1, cycler.Get())
}

func (s *CyclerTestSuite) TestNext() {
    values := []interface{}{1, 2, 3}
    cycler := &cycler{
        values: values,
    }
    s.Require().Equal(2, cycler.Next())
    s.Require().Equal(3, cycler.Next())
    s.Require().Equal(1, cycler.Next())
}

func (s *CyclerTestSuite) TestSeek() {
    values := []interface{}{1, 2, 3}
    cycler := &cycler{
        values: values,
    }
    target := 3
    curValue := cycler.Get()
    value := cycler.Seek(target)
    s.Require().NotEqual(curValue, value)
    s.Require().Equal(3, cycler.Get())
}

func (s *CyclerTestSuite) TestSeekSamePointers() {
    v1 := 1
    v2 := 2
    v3 := 3
    values := []interface{}{&v1, &v2, &v3}
    cycler := &cycler{
        values: values,
    }
    target := &v3
    curValue := cycler.Get()
    value := cycler.Seek(&target)
    s.Require().Equal(curValue, value)
    s.Require().NotEqual(v3, cycler.Get())
}

func (s *CyclerTestSuite) TestSeekPointersFail() {
    v1 := 1
    v2 := 2
    v3 := 3
    values := []interface{}{&v1, &v2, &v3}
    cycler := &cycler{
        values: values,
    }
    target := 3
    curValue := cycler.Get()
    value := cycler.Seek(&target)
    s.Require().Equal(curValue, value)
    s.Require().NotEqual(v3, cycler.Get())
}

func (s *CyclerTestSuite) TestSeekStruct() {
    type V struct {
        Value int
    }
    type T struct {
        Value int
        SubT V
    }
    t1 := T{ Value: 1, SubT: V{ Value: 1 } }
    t2 := T{ Value: 2, SubT: V{ Value: 2 } }
    t3 := T{ Value: 3, SubT: V{ Value: 3 } }
    values := []interface{}{t1, t2, t3}
    cycler := &cycler{
        values: values,
    }
    target := T{ Value: 3, SubT: V{ Value: 3 } }
    curValue := cycler.Get()
    value := cycler.Seek(target)
    s.Require().NotEqual(curValue, value)
    s.Require().Equal(t3, cycler.Get())
}

func (s *CyclerTestSuite) TestSeekNotFound() {
    values := []interface{}{1, 2, 3}
    cycler := &cycler{
        values: values,
    }
    curValue := cycler.Get()
    value := cycler.Seek(4)
    s.Require().Equal(curValue, value)
    s.Require().Equal(curValue, cycler.Get())
}

func (s *CyclerTestSuite) TestSeekNil() {
    values := []interface{}{1, 2, 3}
    cycler := &cycler{
        values: values,
    }
    curValue := cycler.Get()
    value := cycler.Seek(nil)
    s.Require().Equal(curValue, value)
    s.Require().Equal(curValue, cycler.Get())
}

func (s *CyclerTestSuite) TestLen() {
    values := []interface{}{1, 2, 3}
    cycler := &cycler{
        values: values,
    }
    s.Require().Equal(len(values), cycler.Len())
}

func (s *CyclerTestSuite) TestNewCyclerFromSlice() {
    values := []int{1, 2, 3}
    iface, err := NewCyclerFromSlice(values)
    cycler := iface.(*cycler)
    s.Require().NoError(err)
    s.Require().Len(cycler.values, 3)
    for i := range values {
        s.Require().Equal(values[i], cycler.values[i].(int))
    }
}

func (s *CyclerTestSuite) TestNewCyclerFromSliceEmptySlice() {
    values := make([]int, 0)
    _, err := NewCyclerFromSlice(values)
    s.Require().Equal(ErrNoValues, err)
}

func (s *CyclerTestSuite) TestNewCyclerFromSliceNotSlice() {
    values := 1
    _, err := NewCyclerFromSlice(values)
    s.Require().Equal(err, ErrNotSlice)
}

func (s *CyclerTestSuite) TestNewCycler() {
    values := []int{1, 2, 3}
    iface, err := NewCycler(1, 2, 3)
    cycler := iface.(*cycler)
    s.Require().NoError(err)
    s.Require().Len(cycler.values, 3)
    for i := range values {
        s.Require().Equal(values[i], cycler.values[i])
    }
}

func (s *CyclerTestSuite) TestNewCyclerEmpty() {
    _, err := NewCycler()
    s.Require().Equal(err, ErrNoValues)
}
