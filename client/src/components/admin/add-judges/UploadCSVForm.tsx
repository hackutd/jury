import { useState } from 'react';
import { createHeaders } from '../../../api';
import Loading from '../../Loading';
import { useAdminStore } from '../../../store';
import Button from '../../Button';
import Card from '../../Card';

interface UploadCSVFormProps {
    /* The format of the CSV file */
    format: 'project' | 'judge' | 'devpost';
}

const UploadCSVForm = (props: UploadCSVFormProps) => {
    const [file, setFile] = useState<File | null>();
    const [fileName, setFileName] = useState('No file chosen');
    const [headerRow, setHeaderRow] = useState(false);
    const [error, setError] = useState<string | null>();
    const [msg, setMsg] = useState<string | null>();
    const [isUploading, setIsUploading] = useState<boolean>(false);
    const fetchJudgeStats = useAdminStore((state) => state.fetchJudgeStats);
    const fetchProjectStats = useAdminStore((state) => state.fetchProjectStats);

    // Handle file drag and drop
    const handleDrop: React.DragEventHandler<HTMLDivElement> = (e) => {
        e.preventDefault();

        // Make sure that the drag event has a file
        if (e.dataTransfer.items) {
            // Use the first file
            const file = e.dataTransfer.items[0].getAsFile();
            if (file) {
                setFile(file);
                setFileName(file.name);
                setError(null);
                setMsg(null);
            }
        }
    };

    // Handle file upload
    const UploadCSV = async () => {
        setIsUploading(true);
        setError(null);

        // Create the form data object
        const formData = new FormData();
        formData.append('csv', file as Blob);
        formData.append('headerRow', headerRow.toString());

        // Upload the file by calling the upload endpoint
        try {
            // Determine the path based on the format
            const path =
                props.format === 'judge'
                    ? '/judge/csv'
                    : props.format === 'project'
                    ? '/project/csv'
                    : '/project/devpost';

            // Make the request
            const response = await fetch(`${import.meta.env.VITE_JURY_URL}${path}`, {
                method: 'POST',
                body: formData,
                headers: createHeaders('admin', false),
            });

            // Throw error if response is not ok
            if (!response.ok) {
                console.error(`HTTP error! Status: ${response.status}`);
                setError(await response.text());
                setIsUploading(false);
                return;
            }

            // Reset the form and show success message
            setFile(null);
            setFileName('No file chosen');
            const resource = props.format === 'judge' ? 'judge' : 'project';
            setMsg(`Added ${resource}(s) successfully!`);
        } catch (err) {
            console.error(err);
            setError(err as string);
            setIsUploading(false);
        }

        props.format === 'judge' ? fetchJudgeStats() : fetchProjectStats();
        setIsUploading(false);
    };

    const displayText =
        props.format === 'project'
            ? 'name, description, url, "Try It" link, video link, and a comma separated challenge list (in quotes)'
            : 'name, email, track (optional), and notes (optional)';

    return (
        <Card>
            <div className="flex flex-col items-start h-full">
                <h1 className="text-3xl">
                    Upload {props.format === 'devpost' ? 'Devpost CSV' : 'CSV'}
                </h1>
                <p className="text-lg text-light">
                    {props.format === 'devpost'
                        ? 'Upload a CSV file exported from Devpost. Make sure you select Projects data and "do not include personal info".'
                        : `CSV should be formatted in the order of ${displayText} separated by commas; each entry should be on a new line.`}
                </p>
                <div className="flex flex-col w-full space-y-4 mt-4">
                    <div
                        className="flex flex-col items-center justify-center w-full"
                        onDrop={handleDrop}
                        onDragOver={(e) => e.preventDefault()}
                    >
                        <label
                            htmlFor={'dropzone-file-' + props.format}
                            className={`flex flex-col items-center justify-center w-full border-2 border-dashed rounded-sm cursor-pointer 
                                ${
                                    error
                                        ? 'border-error bg-error/20'
                                        : 'border-primary bg-primary/20'
                                }
                                `}
                        >
                            <div
                                className="flex flex-col items-center justify-center pt-5"
                                onDrop={handleDrop}
                                onDragOver={(e) => e.preventDefault()}
                            >
                                <p className="text-2xl text-center">
                                    Drag & Drop files here to upload, <br />
                                    or <span className="text-primary"> Browse Files</span>
                                </p>
                            </div>
                            <input
                                id={'dropzone-file-' + props.format}
                                type="file"
                                className="hidden"
                                onChange={(e) => {
                                    if (e.target.files) {
                                        setFile(e.target.files[0]);
                                        setFileName(e.target.files[0].name);
                                        setError(null);
                                        setMsg(null);

                                        console.log(props.format);
                                    }
                                }}
                            />
                            <div className="flex flex-row items-center justify-center pt-5 pb-6">
                                <input
                                    type="checkbox"
                                    checked={headerRow}
                                    onChange={(e) => {
                                        setHeaderRow(e.target.checked);
                                    }}
                                    className="mr-4 rounded-sm bg-white border-primary border-2 text-primary focus:ring-0"
                                />
                                <p className="text-base text-center">CSV contains a header row</p>
                            </div>
                        </label>
                    </div>
                    {error && (
                        <div className="text-base text-center text-error">
                            Error Uploading CSV: {error}
                        </div>
                    )}
                    {msg && <div className="text-base text-center text-primary">{msg}</div>}
                    {file && (
                        <div className="flex w-full h-11 px-4 text-2xl border-2 border-lightest rounded-md text-start items-center">
                            File Chosen: {fileName}
                        </div>
                    )}
                    <div className="flex flex-row w-full mt-4 space-x-6">
                        <Button
                            type="primary"
                            onClick={UploadCSV}
                            full
                            flat
                            disabled={isUploading || !file || !fileName}
                            className="py-1 rounded-md"
                        >
                            {isUploading ? 'Uploading...' : 'Upload'}
                        </Button>
                    </div>
                </div>
            </div>
            <Loading disabled={!isUploading} />
        </Card>
    );
};

export default UploadCSVForm;
