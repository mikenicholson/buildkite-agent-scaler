package scaler

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type AutoscaleGroupDetails struct {
	DesiredCount int64
	MinSize      int64
	MaxSize      int64
}

type asgDriver struct {
	name string
}

func (a *asgDriver) Describe() (AutoscaleGroupDetails, error) {
	svc := autoscaling.New(session.New())
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(a.name),
		},
	}

	result, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		return AutoscaleGroupDetails{}, err
	}

	return AutoscaleGroupDetails{
		DesiredCount: int64(*result.AutoScalingGroups[0].DesiredCapacity),
		MinSize:      int64(*result.AutoScalingGroups[0].MinSize),
		MaxSize:      int64(*result.AutoScalingGroups[0].MaxSize),
	}, nil

}

func (a *asgDriver) SetDesiredCapacity(count int64) error {
	svc := autoscaling.New(session.New())
	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(a.name),
		DesiredCapacity:      aws.Int64(count),
		HonorCooldown:        aws.Bool(false),
	}

	_, err := svc.SetDesiredCapacity(input)
	if err != nil {
		return err
	}

	return nil
}

type dryRunASG struct {
}

func (a *dryRunASG) Describe() (AutoscaleGroupDetails, error) {
	return AutoscaleGroupDetails{}, nil
}

func (a *dryRunASG) SetDesiredCapacity(count int64) error {
	return nil
}
