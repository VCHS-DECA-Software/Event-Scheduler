package links

import (
	"errors"
	"main/components/common"
	uuid "main/vendor/github.com/satori/go.uuid"
)

type Link struct {
	ID         string `storm:"id"`
	Associated []string
	Maximum    int
}

func NewLink(max int) Link {
	id := uuid.NewV4().String()
	return Link{
		ID:         id,
		Associated: make([]string, 0),
		Maximum:    max,
	}
}

func (self *Link) Add(id string) error {
	if len(self.Associated) >= self.Maximum && self.Maximum > 0 {
		return errors.New("Link has reached maximum capacity.")
	}
	self.Associated = append(self.Associated, id)
	return nil
}

func (self *Link) Has(id string) bool {
	return common.HasElement(id, self.Associated)
}

func (self *Link) Remove(id string) {
	self.Associated = common.RemoveElement(id, self.Associated)
}
