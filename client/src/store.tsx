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
}

const useAdminStore = create<AdminStore>()((set) => ({
    stats: {
        projects: 0,
        avg_project_seen: 0,
        avg_judge_seen: 0,
        judges: 0,
    },

    fetchStats: async () => {
        const statsRes = await getRequest<Stats>('/admin/stats', 'admin');
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
    }
}));

export { useAdminStore, useClockStore };
