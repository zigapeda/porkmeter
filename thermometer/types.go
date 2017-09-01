package main

type ThermometerConfig struct {
	Meters	[]Meter  `yaml:"meters"`
	Gmin	int		`yaml:"gmin"`
	Gmax 	int		`yaml:"gmax"`
	F1		int		`yaml:"f1"`
	F2		int		`yaml:"f2"`
	T1		float64	`yaml:"t1"`
	T2		float64	`yaml:"t2"`
}

type Meter struct {
	Type	string	`yaml:"type"`
	PhysID	string	`yaml:"physid"`
	Channel int		`yaml:"channel"`
	R1		float64 `yaml:"r1"`
	R2		float64	`yaml:"r2"`
	B		float64 `yaml:"b"`
	Temp    int     `yaml:"-"`
	Reading bool	`yaml:"-"`
}