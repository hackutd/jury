import JudgesTable from "./JudgesTable";
import ProjectsTable from "./ProjectsTable";

const AdminTable = (props: { showProjects: boolean }) => {
    return props.showProjects ? <ProjectsTable /> : <JudgesTable />;
};

export default AdminTable;
