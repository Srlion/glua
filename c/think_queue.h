#pragma once
#include <stdio.h>

extern void taskQueueThink(uintptr_t L);

extern void increment_tasks_count();
extern void decrement_tasks_count_by(int n);
extern void reset_tasks_count();
extern int think_queue_think(uintptr_t L);
