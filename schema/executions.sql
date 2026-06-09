-- Async LLM execution tracking.
-- Created by the tasks-producer (llmCall) as 'pending', updated by the
-- queue-handler when the LLM finishes, and polled by the workflow (getResult).
CREATE TABLE executions (
    execution_id Uuid NOT NULL,
    player_id    Uuid,
    chat_id      Uuid,
    status       Utf8,          -- 'pending' | 'done' | 'error'
    response     Utf8,
    error        Utf8,
    create_time  Datetime,
    update_time  Datetime,
    PRIMARY KEY (execution_id)
);
