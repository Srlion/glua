#pragma once
#include <stdio.h>
#include "think_queue.h"
#include "glua.h"

extern void thinkQueueProcess(lua_State L);

static unsigned int tasks_count = 0;

void increment_tasks_count()
{
    tasks_count++;
}

void decrement_tasks_count_by(int n)
{
    tasks_count -= n;
}

void reset_tasks_count()
{
    tasks_count = 0;
}

int think_queue_think(lua_State L)
{
    if (tasks_count > 0)
    {
        thinkQueueProcess(L);
    }
    return 0;
}
