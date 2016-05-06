package common

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
def create_gen_id():
    """根据时间生成一个id """
    gen_id = str(datetime.datetime.now()).replace(' ', '').replace('-','').replace(':', '').replace('.', '')
    gen_id += str(random.randint(0,9))
    return gen_id


def get_uuid():
    """生成一个唯一的用户id
    """
    return md5.md5(create_gen_id() + UUID_STR).hexdigest()
*/

const (
	UUID_STR = "fWv3wFvwSIJ0RuNthkCBeRXnfk5635kufiD5G84MCRcDydAmpD0zxE8QsPjGsAsI"
	UPWD_STR = "JpHot0lcJascWXF5lGP5YNTiKvEEf6RUfrCQI95R7QMIEJQej73CaWMNFgsrm0Ho"
)

func GetUUID() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(CreateGenID()+UUID_STR)))
}

func CreateGenID() string {
	replacer := strings.NewReplacer(" ", "", "-", "", ":", "", ".", "", "+", "")
	gen_id := replacer.Replace(time.Now().String())
	gen_id += fmt.Sprint(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
	return gen_id
}
