import { useEffect, useState } from 'react';
import { fileToString, pad, parseCSV } from '../../util';

interface CSVPreviewProps {
    /* The data of the CSV form */
    data: CSVFormState;
}

const CSVPreview = (props: CSVPreviewProps) => {
    const [csvData, setCsvData] = useState<string[][]>([]);

    useEffect(() => {
        if (!props.data.file) return;

        async function loadData() {
            // TODO: We assume the delimiter is a comma, but we should allow the user to change it
            const rawData = await fileToString(props.data.file as File);
            setCsvData(parseCSV(rawData, ','));
            console.log(parseCSV(rawData, ','));
        }

        loadData();
    }, [props.data.file]);

    if (!props.data.file) return null;

    return (
        <div className="w-full flex flex-col p-2 border-lightest border-2 rounded-md">
            <h2 className="text-center text-2xl font-bold text-light">CSV Preview</h2>
            <div className="w-full overflow-x-auto max-h-[500px]">
                <table className="table-fixed w-full text-lg">
                    <tbody>
                        {csvData.map((row, index) => (
                            <tr key={index} className="border-b border-lightest border-spacing-4">
                                {row.map((cell, cellIndex) => (
                                    <td
                                        key={cellIndex}
                                        title={cell}
                                        className="px-4 py-2 border-b border-lightest text-left w-[150px] font-mono text-sm"
                                    >
                                        {pad(cell, 10)}
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default CSVPreview;
