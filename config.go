package main

import (
	"os"
	"time"
)

type Configuration struct {
	General struct {
		RunStartTime            time.Time
		GitVersion              string // Hash of latest commit
		HostName                string // The machine this was run on - the version number will likely only exist on this machine
		DeckFileName            string `yaml:"deck file name"`
		Decks                   string `yaml:"decks"`                        // must be "consecutive" or "list"
		FirstDeckNum            int    `yaml:"first deck number"`            // must be non-negative integer
		NumberOfDecksToBePlayed int    `yaml:"number of decks to be played"` //must be non-negative integer
		List                    []int  // `yaml:"list"`
		TypeOfPlay              string `yaml:"type of play"` // must be "playOrig" or "playAll"
		Verbose                 int    `yaml:"verbose"`
		OutputTo                string `yaml:"outputTo"`
		outWriter               *os.File
	} `yaml:"general"`
	PlayOrig struct {
		Length          int `yaml:"length of initial override strategy"`
		GameLengthLimit int `yaml:"game length limit in moves made"`
	} `yaml:"play original"`
	PlayAll struct {
		GameLengthLimit int `yaml:"game length limit in million moves tried"`
		//FindAllWinStrats bool `yaml:"find all winning strategies?"`
		ReportingType struct {
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
		PrintWinningMoves   bool   `yaml:"print winning moves"`
		ProgressCounter     int    `yaml:"progress counter in millions"`
		WinLossReport       bool   `yaml:"print final deck by deck win loss record"`
		SaveResultsToSQL    bool   `yaml:"save results to SQL"`
		SQLConnectionString string `yaml:"sql connection string"`
	} `yaml:"play all moves"`
}

/*type Configuration2 struct {
	General  General  `yaml:"general"`
	PlayOrig PlayOrig `yaml:"play original"`
	PlayAll  PlayAll  `yaml:"play all moves"`
}

type ReportingType struct {
	DeckByDeck  bool `yaml:"deck by deck"` // referred to as "DbD_R", "DbD_S" or "DbD_VS", in calls to prntMDet and calls thereto
	MoveByMove  bool `yaml:"move by move"` // referred to as "MbM_R", "MbM_S" or "MbM_VS", in calls to prntMDet and calls thereto
	Tree        bool `yaml:"tree"`         // referred to as "Tree_R", "Tree_N" or "Tree_VN", in calls to prntMDet and calls thereto
	NoReporting bool //not part of yaml file, derived after yaml file is unmarshalled & validated   CONSIDER DELETING
}

type General struct {
	RunStartTime            time.Time
	GitVersion              string // Hash of latest commit
	HostName                string // The machine this was run on - the version number will likely only exist on this machine
	DeckFileName            string `yaml:"deck file name"`
	Decks                   string `yaml:"decks"`                        // must be "consecutive" or "list"
	FirstDeckNum            int    `yaml:"first deck number"`            // must be non-negative integer
	NumberOfDecksToBePlayed int    `yaml:"number of decks to be played"` //must be non-negative integer
	List                    []int  // `yaml:"list"`
	TypeOfPlay              string `yaml:"type of play"` // must be "playOrig" or "playAll"
	Verbose                 int    `yaml:"verbose"`
	OutputTo                string `yaml:"outputTo"`
	outWriter               *os.File
}

type PlayOrig struct {
	Length          int `yaml:"length of initial override strategy"`
	GameLengthLimit int `yaml:"game length limit in moves made"`
}

type DeckByDeckReportingOptions struct {
	Type string `yaml:"type"`
}

type MoveByMoveReportingOptions struct {
	Type string `yaml:"type"`
}

type TreeReportingOptions struct {
	Type                        string `yaml:"type"`
	TreeSleepBetwnMoves         int    `yaml:"sleep between moves"`
	TreeSleepBetwnMovesDur      time.Duration
	TreeSleepBetwnStrategies    int `yaml:"sleep between strategies"`
	TreeSleepBetwnStrategiesDur time.Duration
}

type RestrictReportingTo struct {
	DeckStartVal          int `yaml:"starting deck number"`
	DeckContinueFor       int `yaml:"continue for how many decks"`
	MovesTriedStartVal    int `yaml:"starting move number"`
	MovesTriedContinueFor int `yaml:"continue for how many moves"`
}

type PlayAll struct {
	GameLengthLimit              int                        `yaml:"game length limit in million moves tried"`
	FindAllWinStrats             bool                       `yaml:"find all winning strategies?"`
	ReportingType                ReportingType              `yaml:"reporting"`
	DeckByDeckReportingOptions   DeckByDeckReportingOptions `yaml:"deck by deck reporting options"`
	MoveByMoveReportingOptions   MoveByMoveReportingOptions `yaml:"move by move reporting options"`
	TreeReportingOptions         TreeReportingOptions       `yaml:"tree reporting options"`
	RestrictReporting            bool                       //not part of yaml file, derived after yaml file is unmarshalled & validated
	RestrictReportingTo          RestrictReportingTo        `yaml:"restrict reporting to"`
	PrintWinningMoves            bool                       `yaml:"print winning moves"`
	ProgressCounter              int                        `yaml:"progress counter in millions"`
	WinLossReport                bool   `yaml:"print final deck by deck win loss record"`
	SaveResultsToSQL             bool   `yaml:"save results to SQL"`
	SQLConnectionString          string `yaml:"sql connection string"`
}

type ConfigurationSubsetForSQLWriting struct { // STAN not sure we even need to create this it is simply here for me to communicatewhat needs to be written
	RunStartTime time.Time
	GitVersion   string // Stan I figured out how to do this and will write it tomorrow
	GitSystem    string // The machine this was run on - the version number will likely only exist on this machine
	General      struct {
		Verbose  int    `yaml:"verbose"`
		OutputTo string `yaml:"outputTo"`
	}
	PlayAll struct {
		GameLengthLimit  int  `yaml:"game length limit in million moves tried"`
		FindAllWinStrats bool `yaml:"find all winning strategies?"`
		ReportingType    struct {
			DeckByDeck bool `yaml:"deck by deck"` // referred to as "DbD_R", "DbD_S" or "DbD_VS", in calls to prntMDet and calls thereto
			MoveByMove bool `yaml:"move by move"` // referred to as "MbM_R", "MbM_S" or "MbM_VS", in calls to prntMDet and calls thereto
			Tree       bool `yaml:"tree"`         // referred to as "Tree_R", "Tree_N" or "Tree_VN", in calls to prntMDet and calls thereto
		} `yaml:"reporting"`
		DeckByDeckReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"deck by deck reporting options"`
		MoveByMoveReportingOptions struct {
			Type string `yaml:"type"`
		} `yaml:"move by move reporting options"`
		TreeReportingOptions struct {
			Type                     string `yaml:"type"`
			TreeSleepBetwnMoves      int    `yaml:"sleep between moves"`
			TreeSleepBetwnStrategies int    `yaml:"sleep between strategies"`
		} `yaml:"tree reporting options"`
		RestrictReporting   bool //not part of yaml file, derived after yaml file is unmarshalled & validated
		RestrictReportingTo struct {
			DeckStartVal          int `yaml:"starting deck number"`
			DeckContinueFor       int `yaml:"continue for how many decks"`
			MovesTriedStartVal    int `yaml:"starting move number"`
			MovesTriedContinueFor int `yaml:"continue for how many moves"`
		} `yaml:"restrict reporting to"`
		PrintWinningMoves bool `yaml:"print winning moves"`
		ProgressCounter   int  `yaml:"progress counter in millions"`
	} `yaml:"play all moves"`
}
*/
