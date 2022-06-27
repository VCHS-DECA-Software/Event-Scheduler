package main

import "Event-Scheduler/components/proto"

var judges = []*proto.Judge{
	{
		Firstname: "Lanlan",
		Lastname:  "Gu",
		Judgeable: []string{"PFN"},
	},
	{
		Firstname: "Hui",
		Lastname:  "Cheng",
		Judgeable: []string{"PBM"},
	},
	{
		Firstname: "Joyce",
		Lastname:  "Cheng",
		Judgeable: []string{"PBM", "HRM"},
	},
	{
		Firstname: "Carine",
		Lastname:  "Fang",
		Judgeable: []string{"PBM"},
	},
	{
		Firstname: "Mike",
		Lastname:  "Housholder",
		Judgeable: []string{"PBM"},
	},
	{
		Firstname: "Yuhong",
		Lastname:  "Chen",
		Judgeable: []string{"PMK", "PBM"},
	},
	{
		Firstname: "Manoj",
		Lastname:  "Jain",
		Judgeable: []string{"PHT"},
	},
	{
		Firstname: "Rena",
		Lastname:  "Jin",
		Judgeable: []string{"PHT"},
	},
	{
		Firstname: "Anthony",
		Lastname:  "Lam",
		Judgeable: []string{"PHT", "FS"},
	},
	{
		Firstname: "Sylvia",
		Lastname:  "Lim",
		Judgeable: []string{"PMK"},
	},
	{
		Firstname: "Yan",
		Lastname:  "Liu",
		Judgeable: []string{"PMK"},
	},
	{
		Firstname: "Jennifer",
		Lastname:  "Liu",
		Judgeable: []string{"PMK", "RM"},
	},
	{
		Firstname: "Zhongsheng",
		Lastname:  "Liu",
		Judgeable: []string{"HR"},
	},
	{
		Firstname: "Ritesh",
		Lastname:  "Patel",
		Judgeable: []string{"ENT", "FS"},
	},
	{
		Firstname: "Ling",
		Lastname:  "Wang",
		Judgeable: []string{"MS"},
	},
	{
		Firstname: "Lu",
		Lastname:  "Yang",
		Judgeable: []string{"MS"},
	},
	{
		Firstname: "Michael",
		Lastname:  "Liu",
		Judgeable: []string{"BOR"},
	},
	{
		Firstname: "Tracy",
		Lastname:  "Jiang",
		Judgeable: []string{"BOR"},
	},
	{
		Firstname: "Kerry",
		Lastname:  "Kirchenbauer",
		Judgeable: []string{"BOR", "IM"},
	},
	{
		Firstname: "Yona",
		Lastname:  "Lu",
		Judgeable: []string{"HSW", "PSE"},
	},
	{
		Firstname: "Lisa",
		Lastname:  "Ramirez",
		Judgeable: []string{"EW"},
	},
	{
		Firstname: "Minakshi",
		Lastname:  "Roychoudhury",
		Judgeable: []string{"EW"},
	},
	{
		Firstname: "Melanie",
		Lastname:  "Volpicella",
		Judgeable: []string{"IM", "EW"},
	},
	{
		Firstname: "Julie",
		Lastname:  "Yeomans",
		Judgeable: []string{"IM"},
	},
	{
		Firstname: "Ada",
		Lastname:  "Yue",
		Judgeable: []string{"IM"},
	},
}

var events = []*proto.Event{
	{Id: "PBM"},
	{Id: "PFN"},
	{Id: "PHT"},
	{Id: "PMK"},

	{Id: "BLTDM"},
	{Id: "BTDM"},
	{Id: "ETDM"},
	{Id: "FTDM"},
	{Id: "HTDM"},
	{Id: "MTDM"},
	{Id: "STDM"},
	{Id: "TTDM"},

	{Id: "PFL"},

	{Id: "ACT"},
	{Id: "AAM"},
	{Id: "ASM"},
	{Id: "BFS"},
	{Id: "BSM"},
	{Id: "ENT"},
	{Id: "FMS"},
	{Id: "HLM"},
	{Id: "HRM"},
	{Id: "MCS"},
	{Id: "QSRM"},
	{Id: "RFSM"},
	{Id: "RMS"},
	{Id: "SEM"},

	{Id: "PMBS"},
	{Id: "PMCD"},
	{Id: "PMCA"},
	{Id: "PMCG"},
	{Id: "PMFL"},
	{Id: "PMSP"},

	{Id: "EIP"},
	{Id: "ESB"},
	{Id: "EIB"},
	{Id: "IBP"},
	{Id: "EBG"},
	{Id: "EFB"},

	{Id: "IMCE"},
	{Id: "IMCP"},
	{Id: "IMCS"},

	{Id: "FCE"},
	{Id: "HTPS"},
	{Id: "PSE"},

	{Id: "BOR"},
	{Id: "BMOR"},
	{Id: "HTOR"},
	{Id: "FOR"},
	{Id: "SEOR"},
}
