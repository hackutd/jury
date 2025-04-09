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
    track_seen: { [track: string]: number };
    active: boolean;
    prioritized: boolean;
    score: number;
    stars: number;
    track_stars: { [track: string]: number };
    group: number;
    last_activity: number;
}

interface PublicProject {
    name: string;
    location: number;
    group: number;
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
    track: string;
    notes: string;
    read_welcome: boolean;
    seen: number;
    group_seen: number;
    seen_projects: JudgedProject[];
    rankings: string[];
    active: boolean;
    group: number;
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
    notes: string;
    starred: boolean;
    name: string;
    location: number;
    description: string;
}

type JudgedProjectWithUrl = {
    url: string;
} & JudgedProject;

type SortableJudgedProject = {
    id: number;
} & JudgedProject;

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
    resolved: boolean;
}

interface Options {
    curr_table_num: number;
    clock: ClockState;
    judging_timer: number;
    min_views: number;
    clock_sync: boolean;
    judge_tracks: boolean;
    tracks: string[];
    multi_group: boolean;
    num_groups: number;
    group_sizes: number[];
    group_table_nums: number[][];
    switching_mode: string;
    auto_switch_prop: number;
    manual_switches: number;
    deliberation: boolean;
    group_names: string[];
    ignore_tracks: string[];
    block_reqs: boolean;
    max_req_per_min: number;
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

interface ScoredItem {
    id: string;
    score: number;
}

interface JudgeStats {
    num: number;
    avg_seen: number;
    num_active: number;
}

interface ProjectStats {
    num: number;
    avg_votes: number;
    avg_seen: number;
}

interface Log {
    log: string;
}

interface Code {
    qr_code: string;
}
