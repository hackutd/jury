import JudgeTable from "./JudgeTable";
import ProjectsTable from "./ProjectsTable";

const AdminTable = (props: { showProjects: boolean }) => {
    return props.showProjects ? <ProjectsTable /> : <JudgeTable />;
};

export default AdminTable;
