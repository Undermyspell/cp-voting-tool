export interface VotingEvent {
    EventType: "new_question" | 
               "upvote_question" | 
               "undo_upvote_question" |
               "update_question" |
               "delete_question" | 
               "answer_question" |
               "start_session" |
               "user_connected" |
               "user_disconnected" |
               "stop_session",
    /**
     * JSON stringified
     */
    Payload: string
}