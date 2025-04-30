import { useEffect, useState } from 'react';
import { useAdminStore, useAdminTableStore, useOptionsStore } from '../../../store';
import HeaderEntry from './HeaderEntry';
import { JudgeSortField } from '../../../enums';
import JudgeRow from './JudgeRow';
import AdminTable from './AdminTable';

const JudgesTable = () => {
    const unsortedJudges = useAdminStore((state) => state.judges);
    const fetchJudges = useAdminStore((state) => state.fetchJudges);
    const [sortState, setSortState] = useState<SortState<JudgeSortField>>({
        field: JudgeSortField.None,
        ascending: true,
    });
    const options = useOptionsStore((state) => state.options);
    const selectedTrack = useOptionsStore((state) => state.selectedTrack);
    const judges = useAdminTableStore((state) => state.judges);
    const setJudges = useAdminTableStore((state) => state.setJudges);
    const setSelected = useAdminTableStore((state) => state.setSelected);

    const updateSort = (field: SortField) => {
        if (sortState.field === field) {
            // If sorted on same field and descending, reset sort state
            if (!sortState.ascending) {
                setSortState({
                    field: JudgeSortField.None,
                    ascending: true,
                });
                setJudges(unsortedJudges);
                return;
            }

            // Otherwise, sort descending
            setSortState({
                field,
                ascending: false,
            });
        } else {
            // If in different sorted state, sort ascending on new field
            setSortState({
                field: field as JudgeSortField,
                ascending: true,
            });
        }
    };

    // On load, fetch judges
    useEffect(() => {
        fetchJudges();
    }, [fetchJudges]);

    // When judges change, update judges and sort
    useEffect(() => {
        setSelected(Array(unsortedJudges.length).fill(false));

        sortAndFilterJudges();
    }, [sortState, selectedTrack]);

    useEffect(() => {
        sortAndFilterJudges();
    }, [unsortedJudges]);

    const sortAndFilterJudges = () => {
        // Filter by track
        // TODO: lowk this looks like hot garbage
        const filteredJudges = options.judge_tracks
            ? unsortedJudges.filter(
                  (j) => j.track === selectedTrack || (j.track === '' && selectedTrack === '')
              )
            : unsortedJudges;

        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        let sortFunc = (a: Judge, b: Judge) => 0;
        const asc = sortState.ascending ? 1 : -1;
        switch (sortState.field) {
            case JudgeSortField.Name:
                sortFunc = (a, b) => a.name.localeCompare(b.name) * asc;
                break;
            case JudgeSortField.Code:
                sortFunc = (a, b) => a.code.localeCompare(b.code) * asc;
                break;
            case JudgeSortField.Track:
                sortFunc = (a, b) => a.track.localeCompare(b.track) * asc;
                break;
            case JudgeSortField.Group:
                sortFunc = (a, b) => (a.group - b.group) * asc;
                break;
            case JudgeSortField.Seen:
                sortFunc = (a, b) => (a.seen - b.seen) * asc;
                break;
            case JudgeSortField.Updated:
                sortFunc = (a, b) => (a.last_activity - b.last_activity) * asc;
                break;
            case JudgeSortField.Curr:
                sortFunc = (a, b) => a.current.localeCompare(b.current) * asc;
                break;
        }
        setJudges(filteredJudges.toSorted(sortFunc));
    };

    return (
        <AdminTable>
            <tr>
                <th className="w-12"></th>
                <HeaderEntry
                    name="Name"
                    updateSort={updateSort}
                    sortField={JudgeSortField.Name}
                    sortState={sortState}
                    align="left"
                />
                <HeaderEntry
                    name="Code"
                    updateSort={updateSort}
                    sortField={JudgeSortField.Code}
                    sortState={sortState}
                />
                {options.multi_group && selectedTrack === '' && (
                    <HeaderEntry
                        name="Group"
                        updateSort={updateSort}
                        sortField={JudgeSortField.Group}
                        sortState={sortState}
                    />
                )}
                <HeaderEntry
                    name="Seen"
                    updateSort={updateSort}
                    sortField={JudgeSortField.Seen}
                    sortState={sortState}
                />
                <HeaderEntry
                    name="Curr Proj"
                    updateSort={updateSort}
                    sortField={JudgeSortField.Curr}
                    sortState={sortState}
                />
                <HeaderEntry
                    name="Updated"
                    updateSort={updateSort}
                    sortField={JudgeSortField.Updated}
                    sortState={sortState}
                />
                <th className="text-right w-24">Actions</th>
            </tr>
            {judges.map((judge: Judge, idx) => (
                <JudgeRow key={idx} idx={idx} judge={judge} />
            ))}
        </AdminTable>
    );
};

export default JudgesTable;
