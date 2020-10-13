package primitivepersistence

import (
	"context"

	"github.com/paulwrubel/photolum/config"
	"github.com/sirupsen/logrus"
)

type Primitive struct {
	PrimitiveName             string
	PrimitiveType             string
	EncapsulatedPrimitiveName *string
	A                         []float64
	B                         []float64
	C                         []float64
	ANormal                   []float64
	BNormal                   []float64
	CNormal                   []float64
	Point                     []float64
	Normal                    []float64
	Center                    []float64
	Axis                      *string
	Displacement              []float64
	AxisAngles                []float64
	RotationOrder             *string
	Radius                    *float64
	InnerRadius               *float64
	OuterRadius               *float64
	Height                    *float64
	Angle                     *float64
	Density                   *float64
	IsCulled                  *bool
	HasNegativeNormal         *bool
	HasInvertedNormals        *bool
}

var entity = "primitive"

func Save(plData *config.PhotolumData, baseLog *logrus.Entry, primitive *Primitive) error {
	event := "save"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	tag, err := plData.DB.Exec(context.Background(), `
		INSERT INTO primitives (
			primitive_name,
			primitive_type,
			encapsulated_primitive_name,
			a,
			b,
			c,
			a_normal,
			b_normal,
			c_normal,
			point,
			normal,
			center,
			axis,
			displacement,
			axis_angles,
			rotation_order,
			radius,
			inner_radius,
			outer_radius,
			height,
			angle,
			density,
			is_culled,
			has_negative_normal,
			has_inverted_normals
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25)`,
		primitive.PrimitiveName,
		primitive.PrimitiveType,
		primitive.EncapsulatedPrimitiveName,
		primitive.A,
		primitive.B,
		primitive.C,
		primitive.ANormal,
		primitive.BNormal,
		primitive.CNormal,
		primitive.Point,
		primitive.Normal,
		primitive.Center,
		primitive.Axis,
		primitive.Displacement,
		primitive.AxisAngles,
		primitive.RotationOrder,
		primitive.Radius,
		primitive.InnerRadius,
		primitive.OuterRadius,
		primitive.Height,
		primitive.Angle,
		primitive.Density,
		primitive.IsCulled,
		primitive.HasNegativeNormal,
		primitive.HasInvertedNormals,
	)
	if err != nil || tag.RowsAffected() != 1 {
		return err
	}

	log.Trace("database event completed")
	return nil
}

func Get(plData *config.PhotolumData, baseLog *logrus.Entry, primitiveName string) (*Primitive, error) {
	event := "get"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	primitive := &Primitive{}
	err := plData.DB.QueryRow(context.Background(), `
		SELECT 
			primitive_name,
			primitive_type,
			encapsulated_primitive_name,
			a,
			b,
			c,
			a_normal,
			b_normal,
			c_normal,
			point,
			normal,
			center,
			axis,
			displacement,
			axis_angles,
			rotation_order,
			radius,
			inner_radius,
			outer_radius,
			height,
			angle,
			density,
			is_culled,
			has_negative_normal,
			has_inverted_normals
		FROM primitives
		WHERE primitive_name = $1`, primitiveName).Scan(
		&primitive.PrimitiveName,
		&primitive.PrimitiveType,
		&primitive.EncapsulatedPrimitiveName,
		&primitive.A,
		&primitive.B,
		&primitive.C,
		&primitive.ANormal,
		&primitive.BNormal,
		&primitive.CNormal,
		&primitive.Point,
		&primitive.Normal,
		&primitive.Center,
		&primitive.Axis,
		&primitive.Displacement,
		&primitive.AxisAngles,
		&primitive.RotationOrder,
		&primitive.Radius,
		&primitive.InnerRadius,
		&primitive.OuterRadius,
		&primitive.Height,
		&primitive.Angle,
		&primitive.Density,
		&primitive.IsCulled,
		&primitive.HasNegativeNormal,
		&primitive.HasInvertedNormals,
	)
	if err != nil {
		return nil, err
	}

	log.Trace("database event completed")
	return primitive, nil
}

func Update(plData *config.PhotolumData, baseLog *logrus.Entry, primitive *Primitive) error {
	return nil
}

func Delete(plData *config.PhotolumData, baseLog *logrus.Entry, primitive *Primitive) error {
	return nil
}

func DoesExist(plData *config.PhotolumData, baseLog *logrus.Entry, primitiveName string) (bool, error) {
	event := "exist"
	log := baseLog.WithFields(logrus.Fields{
		"entity": entity,
		"event":  event,
	})
	log.Trace("database event initiated")

	var count int
	err := plData.DB.QueryRow(context.Background(), `
		SELECT count(*)
		FROM primitives
		WHERE primitive_name = $1`, primitiveName).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Trace("database event completed")
	return count == 1, nil
}
