#pragma once
#include <stdio.h>
#include <stdatomic.h>
#include "glua.h"

extern void taskQueueThink(lua_State L);

extern void increment_tasks_count();
extern void decrement_tasks_count_by(int n);
extern void reset_tasks_count();
extern int think_queue_think(lua_State L);
