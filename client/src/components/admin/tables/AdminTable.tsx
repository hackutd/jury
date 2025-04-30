import React from 'react';

const AdminTable = ({ children }: { children: React.ReactNode }) => {
    return (
        <div className="lg:w-full overflow-x-auto lg:overflow-x-clip px-8 pb-4">
            <table className="table-auto lg:table-fixed w-full text-lg">
                <tbody>{children}</tbody>
            </table>
        </div>
    );
};

export default AdminTable;
