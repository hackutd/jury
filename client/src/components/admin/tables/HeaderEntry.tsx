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

    /* Text alignment */
    align?: 'left' | 'center' | 'right';
}

const HeaderEntry = <T extends SortField>(props: HeaderEntryProps<T>) => {
    const handleClick = () => {
        props.updateSort(props.sortField);
    };

    return (
        <th
            className={
                'text-left py-1 cursor-pointer hover:text-primary duration-100' +
                (props.sortState.field === props.sortField ? ' text-primary' : ' text-black') +
                (props.align ? ' text-' + props.align : ' text-center')
            }
            onClick={handleClick}
        >
            {props.name}{' '}
            {props.sortState.field === props.sortField && arrow(props.sortState.ascending)}
        </th>
    );
};

export default HeaderEntry;
