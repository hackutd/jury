import { create } from 'zustand';
import { getRequest } from './api';
import { errorAlert } from './util';

interface AdminStore {
    stats: Stats;
    fetchStats: () => Promise<void>;
    projects: Project[];
    fetchProjects: () => Promise<void>;
    judges: Judge[];
    fetchJudges: () => Promise<void>;
    judgeStats: JudgeStats;
    fetchJudgeStats: () => Promise<void>;
    projectStats: ProjectStats;
    fetchProjectStats: () => Promise<void>;
}

const useAdminStore = create<AdminStore>()((set) => ({
    stats: {
        projects: 0,
        avg_project_seen: 0,
        avg_judge_seen: 0,
        judges: 0,
    },

    fetchStats: async () => {
        const selectedTrack = useOptionsStore.getState().selectedTrack;
        const statsRes = await getRequest<Stats>(`/admin/stats/${selectedTrack}`, 'admin');
        if (statsRes.status !== 200) {
            errorAlert(statsRes);
            return;
        }
        set({ stats: statsRes.data as Stats });
    },

    projects: [],

    fetchProjects: async () => {
        const projRes = await getRequest<Project[]>('/project/list', 'admin');
        if (projRes.status !== 200) {
            errorAlert(projRes);
            return;
        }
        const projects = projRes.data as Project[];

        const scoreRes = await getRequest<ScoredItem[]>('/admin/score', 'admin');
        if (scoreRes.status !== 200) {
            errorAlert(scoreRes);
            return;
        }
        const scores = scoreRes.data as ScoredItem[];

        projects.forEach((project) => {
            const score = scores.find((s) => s.id === project.id);
            if (score) {
                project.score = score.score;
            }
        });

        set({ projects: projRes.data as Project[] });
    },

    judges: [],

    fetchJudges: async () => {
        const judgeRes = await getRequest<Judge[]>('/judge/list', 'admin');
        if (judgeRes.status !== 200) {
            errorAlert(judgeRes);
            return;
        }
        set({ judges: judgeRes.data as Judge[] });
    },

    judgeStats: {
        num: 0,
        avg_seen: 0,
        num_active: 0,
    },

    fetchJudgeStats: async () => {
        const statsRes = await getRequest<JudgeStats>('/judge/stats', 'admin');
        if (statsRes.status !== 200) {
            errorAlert(statsRes);
            return;
        }
        set({ judgeStats: statsRes.data as JudgeStats });
    },

    projectStats: {
        num: 0,
        avg_seen: 0,
        avg_votes: 0,
    },

    fetchProjectStats: async () => {
        const statsRes = await getRequest<ProjectStats>('/project/stats', 'admin');
        if (statsRes.status !== 200) {
            errorAlert(statsRes);
            return;
        }
        set({ projectStats: statsRes.data as ProjectStats });
    },
}));

interface ClockStore {
    clock: ClockState;
    fetchClock: () => Promise<void>;
    setTime: (time: number) => void;
}

const useClockStore = create<ClockStore>()((set) => ({
    clock: {
        time: 0,
        running: false,
    },

    fetchClock: async () => {
        const clockRes = await getRequest<ClockState>('/admin/clock', 'admin');
        if (clockRes.status !== 200) {
            errorAlert(clockRes);
            return;
        }
        set({ clock: clockRes.data as ClockState });
    },

    setTime: (time: number) => {
        set((state) => ({
            clock: {
                ...state.clock,
                time: time,
            },
        }));
    },
}));

interface OptionsStore {
    options: Options;
    selectedTrack: string;
    fetchOptions: () => Promise<void>;
    setSelectedTrack: (track: string) => void;
}

const useOptionsStore = create<OptionsStore>((set) => ({
    options: {
        curr_table_num: 0,
        clock: {
            time: 0,
            running: false,
        },
        judging_timer: 0,
        categories: [],
        min_views: 0,
        clock_sync: false,
        judge_tracks: false,
        tracks: [],
        multi_group: false,
        num_groups: 0,
        group_sizes: [],
        group_table_nums: [],
        main_group: {
            switching_mode: '',
            auto_switch_method: '',
            auto_switch_count: 0,
            auto_switch_prop: 0,
            manual_switches: 0,
        },
    },

    selectedTrack: '',

    fetchOptions: async () => {
        const optionsRes = await getRequest<Options>('/admin/options', 'admin');
        if (optionsRes.status !== 200) {
            errorAlert(optionsRes);
            return;
        }
        set({ options: optionsRes.data as Options });
    },

    setSelectedTrack: async (track: string) => {
        // If main judging selected, reset selected track
        if (track === 'Main Judging') {
            set({ selectedTrack: '' });
            return;
        }

        set({ selectedTrack: track });
    },
}));

interface FlagsStore {
    flags: Flag[];
    fetchFlags: () => Promise<void>;
}

const useFlagsStore = create<FlagsStore>((set) => ({
    flags: [],

    fetchFlags: async () => {
        const flagsRes = await getRequest<Flag[]>('/admin/flags', 'admin');
        if (flagsRes.status !== 200) {
            errorAlert(flagsRes);
            return;
        }
        set({ flags: flagsRes.data as Flag[] });
    },
}));

export { useAdminStore, useClockStore, useOptionsStore, useFlagsStore };
