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
        seen: 0,
        votes: 0,
        avg_mu: 0,
        avg_sigma: 0,
        judges: 0,
    },

    fetchStats: async () => {
        const statsRes = await getRequest<Stats>('/admin/stats', 'admin');
        if (statsRes.status !== 200) {
            errorAlert(statsRes.status);
            return;
        }
        set({ stats: statsRes.data as Stats });
    },

    projects: [],

    fetchProjects: async () => {
        const projRes = await getRequest<Project[]>('/project/list', 'admin');
        if (projRes.status !== 200) {
            errorAlert(projRes.status);
            return;
        }
        set({ projects: projRes.data as Project[] });
    },

    judges: [],

    fetchJudges: async () => {
        const judgeRes = await getRequest<Judge[]>('/judge/list', 'admin');
        if (judgeRes.status !== 200) {
            errorAlert(judgeRes.status);
            return;
        }
        set({ judges: judgeRes.data as Judge[] });
    }
}));

export default useAdminStore;
