#pragma once
#include <stdio.h>
#include <stdatomic.h>
#include "think_queue.h"
#include "glua.h"

extern void thinkQueueProcess(lua_State L);

static atomic_uint tasks_count = 0;

void increment_tasks_count()
{
    atomic_fetch_add_explicit(&tasks_count, 1, memory_order_relaxed);
}

void decrement_tasks_count_by(int n)
{
    atomic_fetch_sub_explicit(&tasks_count, n, memory_order_relaxed);
}

void reset_tasks_count()
{
    atomic_store_explicit(&tasks_count, 0, memory_order_relaxed);
}

int think_queue_think(lua_State L)
{
    if (atomic_load_explicit(&tasks_count, memory_order_relaxed) > 0)
    {
        thinkQueueProcess(L);
    }
    return 0;
}
