<?php
interface ChainRootBad { public function dispatch(object $event); }
interface ChainMiddleBad extends ChainRootBad { public function dispatch(object $event, ?string $name = null): object; }
interface ChainLeafBad extends ChainMiddleBad { public function dispatch(object $event): object; }
