package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
)

type LambdaResult struct {
	Message string `json:"Answer:"`
}

func hello() (LambdaResult, error) {
	// 認証
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("認証処理失敗" + err.Error())
	}
	cfg.Region = endpoints.ApNortheast1RegionID

	// 起動していないEC2インスタンス情報を取得
	svc := ec2.New(cfg)
	notStoppedIncetaces, err := GetnotRunningIncetaces(svc)
	if err != nil {
		log.Fatalf("インスタンス情報取得エラー:%v", err)
	}

	// EC2インスタンス起動
	if len(notStoppedIncetaces) == 0 {
		log.Println("停止状態のインスタンスはありません。正常終了します。")
		return LambdaResult{Message: "Success!"}, nil
	}
	result, err := StopIncetances(svc, notStoppedIncetaces)
	if err != nil {
		log.Fatalf("インスタンス起動エラー:%v", err)
	}
	for _, v := range result.StoppingInstances{
		log.Printf("[%v]:[%v] -> [%v]", *v.InstanceId, v.PreviousState.Name, v.CurrentState.Name)
	}
	return LambdaResult{Message: "Success!"}, nil
}

func main() {
	lambda.Start(hello)
}

func GetnotRunningIncetaces(svc *ec2.EC2) ([]string, error) {
	var notRunnigIncetances []string
	params := &ec2.DescribeInstancesInput{}
	request := svc.DescribeInstancesRequest(params)
	result, err := request.Send()
	if err != nil {
		return nil, err
	}
	for insnum, instances := range result.Reservations {
		for _, instance := range instances.Instances {
			log.Printf("Incetance:[%v][%v] State:[%v]\n", insnum, *instance.InstanceId, instance.State.Name)
			if instance.State.Name != "stopped" {
				log.Printf("停止するインスタンスに登録しました[%v]", *instance.InstanceId)
				notRunnigIncetances = append(notRunnigIncetances, *instance.InstanceId)
			}
		}
	}
	return notRunnigIncetances, nil
}

func StopIncetances(svc *ec2.EC2, notRunningIncetaces []string) (*ec2.StopInstancesOutput, error) {
	params := &ec2.StopInstancesInput{
		InstanceIds: notRunningIncetaces,
	}
	request := svc.StopInstancesRequest(params)
	result, err := request.Send()
	if err != nil {
		return nil, err
	}
	return result, nil
}
