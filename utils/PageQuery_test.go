package utils_test

import (
	"gitlab-booster/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PageQuery", func() {
	Describe("EffecitvePageSize and EffecitvePageNum", func() {
		It("shoud generate correct effective value", func() {
			Expect((&utils.PageQuery{}).EffecitvePageNum()).To(Equal(uint(1)))
			Expect((&utils.PageQuery{PageNum: 0}).EffecitvePageNum()).To(Equal(uint(1)))
			Expect((&utils.PageQuery{PageNum: 3}).EffecitvePageNum()).To(Equal(uint(3)))

			Expect((&utils.PageQuery{}).EffecitvePageSize()).To(Equal(uint(10)))
			Expect((&utils.PageQuery{PageSize: 0}).EffecitvePageSize()).To(Equal(uint(10)))
			Expect((&utils.PageQuery{PageSize: 5}).EffecitvePageSize()).To(Equal(uint(5)))

			Expect((&utils.PageQuery{}).Offset()).To(Equal(uint(0)))                         // (1-1) * 10
			Expect((&utils.PageQuery{PageSize: 0, PageNum: 0}).Offset()).To(Equal(uint(0)))  // (1-1) * 10
			Expect((&utils.PageQuery{PageSize: 5, PageNum: 3}).Offset()).To(Equal(uint(10))) // (3-1)*5
		})
	})

	Describe("SortSQL", func() {
		It("should generate corret sort SQL segment", func() {
			Expect((&utils.PageQuery{}).SortSQL()).To(Equal(""))
			Expect((&utils.PageQuery{Sort: "aaa;bbb"}).SortSQL()).To(Equal("`aaa`,`bbb`"))
			Expect((&utils.PageQuery{Sort: "aaa,desC,,;;bbb,Asc;ccc,,"}).SortSQL()).To(Equal("`aaa` DESC,`bbb` ASC,`ccc`"))
			Expect((&utils.PageQuery{Sort: "aaa,desC;bbb,Asc;ccc", Order: "desc"}).SortSQL()).To(Equal("`aaa` DESC,`bbb` ASC,`ccc`"))
			Expect((&utils.PageQuery{Sort: "aaa,Asc", Order: "desc"}).SortSQL()).To(Equal("`aaa` ASC"))
			Expect((&utils.PageQuery{Sort: ",Asc", Order: "desc"}).SortSQL()).To(Equal(""))
			Expect((&utils.PageQuery{Sort: "aaa,xxx,xx", Order: "desc"}).SortSQL()).To(Equal("`aaa`"))

			Expect((&utils.PageQuery{Sort: "aaa", Order: "desc"}).SortSQL()).To(Equal("`aaa` DESC"))
			Expect((&utils.PageQuery{Sort: "aaa", Order: " xxx"}).SortSQL()).To(Equal("`aaa`"))
		})
	})
})
