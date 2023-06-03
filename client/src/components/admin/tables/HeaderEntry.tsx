import { arrow } from '../../../util';

interface HeaderEntryProps<T extends SortField> {
    /* Name of the field */
    name: string;

    /* Function to sort by this field */
    updateSort: (field: T) => void;

    /* Sort state */
    sortState: SortState<T>;

    /* Sort field to use */
    sortField: T;
}

const HeaderEntry = <T extends SortField>(props: HeaderEntryProps<T>) => {
    const handleClick = () => {
        props.updateSort(props.sortField);
    };

    return (
        <th
            className={
                'text-left py-1 cursor-pointer hover:text-primary duration-100' +
                (props.sortState.field === props.sortField ? ' text-primary' : ' text-black')
            }
            onClick={handleClick}
        >
            {props.name}{' '}
            {props.sortState.field === ProjectSortField.Name && arrow(props.sortState.ascending)}
        </th>
    );
};

export default HeaderEntry;
