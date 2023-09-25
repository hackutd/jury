interface ParagraphProps {
    /** Text to display */
    text: string;

    /** Style with className */
    className?: string;
}

/**
 * Formats a paragraph of text.
 * This will split the text into its respective paragraphs,
 * adding a slight space between each body section.
 */
const Paragraph = (props: ParagraphProps) => {
    // If there is no text, return "no description"
    if (props.text === '') return <p className={props.className}>[No Description]</p>;

    // Split the pararaphs by the newline character
    const paragraphs = props.text.split('\\n');

    // Create a list of paragraph elements
    const paragraphElements = paragraphs.map((p, i) => (
        <p key={i} className="mb-3 whitespace-pre-line overflow-x-hidden">
            {p}
        </p>
    ));

    // Return list of items in a div
    return <div className={props.className}>{paragraphElements}</div>;
};

export default Paragraph;
