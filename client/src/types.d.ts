interface Project {
    id: string;
    name: string;
    location: number;
    description: string;
    url: string;
    try_link: string;
    video_link: string;
    challenge_list: string[];
    seen: number;
    active: boolean;
    score: number;
    last_activity: number;
}

interface PublicProject {
    name: string;
    location: number;
    description: string;
    url: string;
    try_link: string;
    video_link: string;
    challenge_list: string;
}

interface Judge {
    id: string;
    name: string;
    code: string;
    email: string;
    notes: string;
    read_welcome: boolean;
    seen: number;
    seen_projects: JudgedProject[];
    rankings: string[];
    active: boolean;
    current: string;
    last_activity: number;
}

interface Stats {
    projects: number;
    avg_project_seen: number;
    avg_judge_seen: number;
    judges: number;
}

type SortField = ProjectSortField | JudgeSortField;

interface SortState<T extends SortField> {
    field: T;
    ascending: boolean;
}

// TODO: Change this...
type VotePopupState = 'vote' | 'busy' | 'flag';

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
    categories: { [name: string]: number };
    notes: string;
    name: string;
    location: number;
    description: string;
}

interface ClockState {
    time: number;
    running: boolean;
}

interface ProjectCount {
    count: number;
}

interface Flag {
    id: string;
    judge_id: string;
    project_id: string;
    time: number;
    project_name: string;
    judge_name: string;
    project_location: number;
    reason: string;
}

interface Options {
    curr_table_num: number;
    clock: ClockState;
    judging_timer: number;
    categories: string[];
}

interface FetchResponse<T> {
    status: number;
    error: string;
    data: T | null;
}

interface Timer {
    judging_timer: number;
}

interface NextJudgeProject {
    project_id: string;
}
