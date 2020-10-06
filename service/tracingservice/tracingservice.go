package tracingservice

import (
	"fmt"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

func StartRender(plData *config.PhotolumData, baseLog *logrus.Logger, renderName string) {
	for i := 0; i < 10000000; i++ {
		fmt.Printf("counting to 10,000,000! We're at %d\n", i)
	}
}
