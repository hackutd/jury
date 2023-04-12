import Button from '../Button';

const AdminToolbar = (props: { showProjects: boolean }) => {
    return (
        <div className="flex flex-row px-8">
            <div>
                <Button
                    type="primary"
                    square
                    className="py-2 px-4"
                    href={props.showProjects ? '/admin/add-projects' : '/admin/add-judges'}
                >
                    Add {props.showProjects ? 'Projects' : 'Judges'}
                </Button>
            </div>
        </div>
    );
};

export default AdminToolbar;
