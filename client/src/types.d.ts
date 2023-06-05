interface Project {
    _id: { $oid: string };
    name: string;
    location: number;
    mu: number;
    sigma_sq: number;
    votes: number;
    seen: number;
    last_activity: {
        $date: {
            $numberLong: number;
        };
    };
}

interface Judge {
    _id: { $oid: string };
    name: string;
    email: string;
    notes: string;
    alpha: number;
    votes: number;
    beta: number;
    last_activity: {
        $date: {
            $numberLong: number;
        };
    };
    read_welcome: boolean;
    next: { $oid: string };
    prev: { $oid: string };
}

interface Stats {
    projects: number;
    seen: number;
    votes: number;
    avg_mu: number;
    avg_sigma: number;
    judges: number;
}

type SortField = ProjectSortField | JudgeSortField;

interface SortState<T extends SortField> {
    field: T;
    ascending: boolean;
}
