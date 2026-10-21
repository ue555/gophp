<?php
function needsCallable(callable $callback): string { return 'accepted'; }
function needsClosure(Closure $callback): string { return 'closure'; }
try { echo needsCallable('notDefined'); } catch (TypeError $e) { echo 'type|'; }
try { echo needsClosure('strlen'); } catch (TypeError $e) { echo 'closure-type|'; }
try { Closure::fromCallable('notDefined'); } catch (TypeError $e) { echo 'from-type|'; }
echo (int)is_callable('notDefined'), '|', (int)is_callable('notDefined', true);
