package main

import "time"

type Configuration struct {
	RunStartTime time.Time
	GitVersion   string // Stan we need to figure out how to get this
	General      struct {
		Decks                   string `yaml:"decks"`                        // must be "consecutive" or "list"
		FirstDeckNum            int    `yaml:"first deck number"`            // must be non-negative integer
		NumberOfDecksToBePlayed int    `yaml:"number of decks to be played"` //must be non-negative integer
		List                    []int
		TypeOfPlay              string `yaml:"type of play"` // must be "playOrig" or "playAll"
		Verbose                 int    `yaml:"verbose"`
		OutputTo                string `yaml:"outputTo"`
	} `yaml:"general"`
	PlayOrig struct {
		Length          int `yaml:"length of initial override strategy"`
		GameLengthLimit int `yaml:"game length limit in moves made"`
	} `yaml:"play original"`
	PlayAll struct {
		GameLengthLimit  int  `yaml:"game length limit in moves tried"`
		FindAllWinStrats bool `yaml:"find all winning strategies?"`
		ReportingType    struct {
			DeckByDeck  bool `yaml:"deck by deck"` // referred to as "DbD_R", "DbD_S" or "DbD_VS", in calls to prntMDet and calls thereto
			MoveByMove  bool `yaml:"move by move"` // referred to as "MbM_R", "MbM_S" or "MbM_VS", in calls to prntMDet and calls thereto
			Tree        bool `yaml:"tree"`         // referred to as "Tree_R", "Tree_N" or "Tree_VN", in calls to prntMDet and calls thereto
			NoReporting bool //not part of yaml file, derived after yaml file is unmarshalled & validated   CONSIDER DELETING
		} `yaml:"reporting"`
		DeckByDeckReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"deck by deck reporting options"`
		MoveByMoveReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"move by move reporting options"`
		TreeReportingOptions struct {
			Type                        string `yaml:"type"`
			TreeSleepBetwnMoves         int    `yaml:"sleep between moves"`
			TreeSleepBetwnMovesDur      time.Duration
			TreeSleepBetwnStrategies    int `yaml:"sleep between strategies"`
			TreeSleepBetwnStrategiesDur time.Duration
		} `yaml:"tree reporting options"`
		RestrictReporting   bool //not part of yaml file, derived after yaml file is unmarshalled & validated
		RestrictReportingTo struct {
			DeckStartVal          int `yaml:"starting deck number"`
			DeckContinueFor       int `yaml:"continue for how many decks"`
			MovesTriedStartVal    int `yaml:"starting move number"`
			MovesTriedContinueFor int `yaml:"continue for how many moves"`
		} `yaml:"restrict reporting to"`
		PrintWinningMoves            bool `yaml:"print winning moves"`
		ProgressCounter              int  `yaml:"progress counter in millions"`
		ProgressCounterLastPrintTime time.Time
		WinLossReport                bool   `yaml:"print final deck by deck win loss record"`
		SaveResultsToSQL             bool   `yaml:"save results to SQL"`
		SQLConnectionString          string `yaml:"sql connection string"`
	} `yaml:"play all moves"`
}
