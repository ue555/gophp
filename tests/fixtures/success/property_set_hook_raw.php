<?php
if (PHP_VERSION_ID < 80400) { echo "hooks unavailable\n"; exit; }
eval(<<<'PHP'
class HookBox {
    public int $calls = 0;
    public int $number = 1 {
        set {
            $this->calls++;
            if ($value < 0) { throw new ValueError('negative'); }
            $this->number = $value * 2;
        }
    }
}
class ChildBox extends HookBox {}
PHP);
$b = new ChildBox;
$b->number = 3;
echo $b->number, '|', $b->calls, "\n";
$r = new ReflectionProperty(HookBox::class, 'number');
$r->setRawValue(object: $b, value: '7');
echo $b->number, '|', $b->calls, '|', $r->getRawValue($b), "\n";
try { $b->number = -1; } catch (ValueError $e) { echo $e->getMessage(), "\n"; }
echo $b->number, '|', $b->calls, "\n";
try { $r->setRawValue($b, []); } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
try { $r->setRawValue(new stdClass, 9); } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
echo $b->number, '|', $b->calls, "\n";
try { $r->getRawValue(); } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
try { $r->setRawValue(unknown: 1); } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
