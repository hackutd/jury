import { useEffect, useState } from 'react';
import { getRequest, postRequest, putRequest } from '../../api';
import { errorAlert } from '../../util';
import Button from '../Button';

interface RatingsProps {
    callback: () => void;
    submitText?: string;
    prior?: { [x: string]: number }; // TODO: wtf is this type
    small?: boolean;
    update?: boolean;
    project: JudgedProject;
}

const Ratings = (props: RatingsProps) => {
    const [categories, setCategories] = useState<string[]>([]);
    const [categoryScores, setCategoryScores] = useState<number[]>([]);

    useEffect(() => {
        if (!props.prior || categories.length === 0) return;

        console.log(props.prior);

        const newScores = categories.map((v) => (props.prior as any)[v] ?? 0); // TODO: fix this
        setCategoryScores(newScores);
    }, [props.prior, categories]);

    // Get categories from the options
    useEffect(() => {
        async function getCategories() {
            // Get the categories
            const res = await getRequest<string[]>(`/categories`, 'judge');
            if (res.status !== 200) {
                errorAlert(res);
                return;
            }

            const cats = res.data ?? [];
            const catScores = cats.map(() => 0);

            setCategories(cats);
            setCategoryScores(catScores);
        }

        getCategories();
    }, []);

    // Submit the scores
    const submit = async () => {
        // Create the scores object
        const scores = categories
            .map((v, i) => ({ [v]: categoryScores[i] }))
            .reduce((a, b) => ({ ...a, ...b }), {});

        // Score the current project
        const scoreRes = props.update
            ? await putRequest<OkResponse>('/judge/score', 'judge', {
                  categories: scores,
                  project: props.project.project_id,
              })
            : await postRequest<OkResponse>('/judge/score', 'judge', {
                  categories: scores,
              });
        if (scoreRes.status !== 200) {
            errorAlert(scoreRes);
            return;
        }

        props.callback();
    };

    return (
        <div className="flex flex-col align-center">
            {categories.map((v, i) => (
                <div key={i}>
                    <p className="text-center">
                        <b>{v}</b>: {categoryScores[i]}
                    </p>
                    <input
                        type="range"
                        min="0"
                        max="10"
                        value={categoryScores[i]}
                        onChange={(e) => {
                            const newScores = [...categoryScores];
                            newScores[i] = parseInt(e.target.value);
                            setCategoryScores(newScores);
                        }}
                        className="w-full"
                    />
                </div>
            ))}
            <div className="flex justify-center mt-4">
                <Button
                    type="primary"
                    onClick={submit}
                    className={props.small ? 'p-1 mb-4' : 'mb-4'}
                >
                    {props.submitText ?? 'Submit'}
                </Button>
            </div>
        </div>
    );
};

export default Ratings;
