package slim4go

/*
import (
	"fmt"
	. "slim4go/slimlist"
	"testing"
)

type Division struct {
	numerator   float32
	denominator float32
	result      string
}

func (self *Division) Quotient(args *SlimList) string {
	quotient := self.numerator / self.denominator
	return strconv.FormatFloat(float64(quotient), 'f', -1, 32)
}

func NewDivision(*SlimList) interface{} {
	return new(Division)
}

func setup(t *testing.T) {
	RegisterFixture("Division", NewDivision)
}

func TestSlimDecisionTable(t *testing.T) {

	message := "[000009:"+
		"000085:[000004:000017:decisionTable_0_0:000004:make:000015:decisionTable_0:000008:Division:]:" +
		"000214:[000005:000017:decisionTable_0_1:000004:call:000015:decisionTable_0:000005:table:000124:" +
			"[000002:" +
				"000062:[000003:000009:numerator:000011:denominator:000009:Quotient?:]:"+
				"000037:[000003:000002:10:000001:5:000001:2:]:]:]:"+
		"000087:[000004:000017:decisionTable_0_2:000004:call:000015:decisionTable_0:000010:beginTable:]:"+
		"000082:[000004:000017:decisionTable_0_3:000004:call:000015:decisionTable_0:000005:reset:]:"+
		"000099:[000005:000017:decisionTable_0_4:000004:call:000015:decisionTable_0:000012:setNumerator:000002:10:]:"+
		"000100:[000005:000017:decisionTable_0_5:000004:call:000015:decisionTable_0:000014:setDenominator:000001:5:]:"+
		"000084:[000004:000017:decisionTable_0_6:000004:call:000015:decisionTable_0:000007:execute:]:"+
		"000085:[000004:000017:decisionTable_0_7:000004:call:000015:decisionTable_0:000008:Quotient:]:"+
		"000085:[000004:000017:decisionTable_0_8:000004:call:000015:decisionTable_0:000008:endTable:]:]"
	response := handleMessage(message)

	fmt.Println(response)
}
*/
