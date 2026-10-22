<?php
if (PHP_VERSION_ID < 80400) { echo "hooks unavailable\n"; exit; }
eval('class HookEdges { public int $x=1 { set { $this->x=$value*2; } } public array $items=[] { set { $this->items=$value; } } }');
$x = new HookEdges;
try { $alias =& $x->x; } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
try { $x->items[] = 1; } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
try { $x->x = []; } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
echo $x->x, '|', count($x->items), "\n";
$r = new ReflectionProperty(HookEdges::class, 'x');
$o = new stdClass;
$r->setRawValue($o, []);
echo json_encode($o), "\n";
try { $r->getRawValue($o); } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
$x->x = 3;
echo $x->x, "\n";
