import JuryHeader from "../../components/JuryHeader";
import UploadCSVForm from "../../components/admin/add-judges/UploadCSVForm";
import AddProjectsStatsPanel from "../../components/admin/add-projects/AddProjectsStatsPanel";
import NewProjectForm from "../../components/admin/add-projects/NewProjectForm";

const AddProjects = () => {
    return (
        <>
            <JuryHeader withBack withLogout isAdmin />
            <div className="flex flex-col items-start justify-center w-full px-8 py-4 md:px-16 md:py-8">
                <h1 className="text-4xl font-bold">Add Projects</h1>
                    <AddProjectsStatsPanel />
                <div className="mt-8 flex flex-col w-full space-y-8">
                    <NewProjectForm />
                    <UploadCSVForm format="project" />
                    <UploadCSVForm format="devpost" />
                </div>
            </div>
        </>
    );
};

export default AddProjects;
