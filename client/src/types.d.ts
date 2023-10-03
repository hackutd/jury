interface Project {
    id: string;
    name: string;
    location: number;
    mu: number;
    sigma_sq: number;
    votes: number;
    seen: number;
    description: string;
    active: boolean;
    prioritized: boolean;
    last_activity: number;
}

interface Judge {
    id: string;
    name: string;
    email: string;
    notes: string;
    alpha: number;
    votes: number;
    beta: number;
    last_activity: number;
    read_welcome: boolean;
    seen_projects: JudgedProject[];
    active: boolean;
    next: string;
    prev: string;
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

interface JudgeVoteRes {
    judge_id: string;
    prev_project_id: string;
    next_project_id: string;
}

type VotePopupState = 'vote' | 'skip' | 'flag';

interface VotingProjectInfo {
    curr_name: string;
    curr_location: number;
    prev_name: string;
    prev_location: number;
}

interface OkResponse {
    ok: number;
}

interface TokenResponse {
    token: string;
}

interface JudgedProject {
    project_id: string;
    name: string;
    description: string;
    stars: number;
}

interface ClockState {
    time: number;
    running: boolean;
}

interface JudgeIpo {
    initial: number;
    project_id: string;
}

interface ProjectCount {
    count: number;
}
