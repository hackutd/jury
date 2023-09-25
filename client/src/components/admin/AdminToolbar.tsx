import Button from '../Button';

const AdminToolbar = (props: { showProjects: boolean }) => {
    return (
        <div className="flex flex-row px-8 py-4">
            <div>
                <Button
                    type="outline"
                    square
                    bold
                    full
                    className="py-2 px-4 rounded-md"
                    href={props.showProjects ? '/admin/add-projects' : '/admin/add-judges'}
                >
                    Add {props.showProjects ? 'Projects' : 'Judges'}
                </Button>
            </div>
        </div>
    );
};

export default AdminToolbar;
