import { create } from 'zustand';

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
        const fetchedStats = await fetch(`${process.env.REACT_APP_JURY_URL}/admin/stats`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        }).then((data) => data.json());
        set({ stats: fetchedStats });
    },

    projects: [],

    fetchProjects: async () => {
        const res = await fetch(`${process.env.REACT_APP_JURY_URL}/project/list`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        });

        if (!res.ok) {
            alert('Error fetching projects');
            return;
        }

        const fetchedProjects = await res.json();
        set({ projects: fetchedProjects });
    },

    judges: [],

    fetchJudges: async () => {
        const res = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/list`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
        });

        if (!res.ok) {
            alert('Error fetching judges');
            return;
        }

        const fetchedJudges = await res.json();
        set({ judges: fetchedJudges });
    }
}));

export default useAdminStore;
