import { useState } from 'react';

const UploadCSVForm = () => {
    const [file, setFile] = useState<File | null>(null);
    const [fileName, setFileName] = useState<string>('No file chosen');
    const [headerRow, setHeaderRow] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);


    const [isUploading, setIsUploading] = useState<boolean>(false);

    const handleDrop: React.DragEventHandler<HTMLDivElement> = (e) => {
        e.preventDefault();
        if (e.dataTransfer.items) {
            const file = e.dataTransfer.items[0].getAsFile();
            if (file) {
                setFile(file);
                setFileName(file.name);
            }
        }
    };

    const UploadCSV = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setIsUploading(true);
        setError(null);

        try {
            const formData = new FormData();
            formData.append('csv', file as Blob);
            formData.append('headerRow', headerRow.toString());
            const response = await fetch(`${process.env.REACT_APP_JURY_URL}/judge/csv/upload`, {
                method: 'POST',
                body: formData,
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            alert('CSV Uploaded');
            setFile(null);
            setFileName('No file chosen');
        } catch (err) {
            console.log(err);
            setError(err as string);
        } finally {
            setIsUploading(false);
        }
    };


    return (
        <>
            <div className="w-full h-full border-lightest border-2 p-8 rounded-sm">
                <div className="flex flex-col items-start h-full">
                    <h1 className="text-3xl">Upload CSV</h1>
                    <p className="text-lg text-light">
                        CSV should be formatted in the order of name, email, and notes separated by
                        commas and then data on newlines
                    </p>
                    <form className="flex flex-col w-full space-y-4 mt-4" onSubmit={UploadCSV}>
                        <div
                            className="flex flex-col items-center justify-center w-full"
                            onDrop={handleDrop}
                            onDragOver={(e) => e.preventDefault()}
                        >
                            <label
                                htmlFor="dropzone-file"
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
                                    id="dropzone-file"
                                    type="file"
                                    className="hidden"
                                    onChange={(e) => {
                                        if (e.target.files) {
                                            setFile(e.target.files[0]);
                                            setFileName(e.target.files[0].name);
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
                                    <p className="text-base text-center">
                                        CSV contains a header row
                                    </p>
                                </div>
                            </label>
                        </div>
                        {error && (
                            <div className="text-base text-center text-error">
                                Error Uploading CSV: {error}
                            </div>
                        )}
                        {file && (
                            <div className="flex w-full h-11 px-4 text-2xl border-2 border-lightest rounded-md text-start items-center">
                                File Chosen: {fileName}
                            </div>
                        )}
                        <div className="flex flex-row w-full mt-4 space-x-6">
                            <button
                                className={`w-full h-11 px-4 text-2xl text-white bg-primary rounded-full
                                ${
                                    isUploading || !file || !fileName
                                        ? 'opacity-50 cursor-not-allowed'
                                        : ''
                                }`}
                                disabled={isUploading || !file || !fileName}
                            >
                                {isUploading ? 'Uploading...' : 'Upload'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </>
    );
};

export default UploadCSVForm;
