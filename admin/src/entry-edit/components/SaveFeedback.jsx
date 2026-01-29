export function SaveFeedback({ feedback }) {
    if (!feedback) return <div class="save-feedback" />;

    const className = feedback.type === 'error' ? 'feedback-error' : 'feedback-success';

    return (
        <div class="save-feedback">
            <div class={className}>{feedback.message}</div>
        </div>
    );
}
