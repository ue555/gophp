<?php
class ConcatBox { public string $text = 'a'; public int $n = 1; }
function stopConcat($level, $message) { throw new Exception('stop'); }
$box = new ConcatBox;
set_error_handler('stopConcat');
try { $box->text .= []; } catch (Exception $e) { echo 'array:', $box->text, '|'; }
try { $box->n .= []; } catch (Exception $e) { echo 'int:', gettype($box->n), $box->n, '|'; }
restore_error_handler();
try { $box->text .= new stdClass; } catch (Error $e) { echo 'object:', $box->text; }
class F { public float $f = 1.5; public ?int $ni = 4; public int|string $u = 3; public mixed $m = 2; }
$g = new F;
set_error_handler('stopConcat');
try { $g->f .= []; } catch (Exception $e) { echo "|f:", gettype($g->f), $g->f; }
try { $g->ni .= []; } catch (Exception $e) { echo "|ni:", gettype($g->ni), $g->ni; }
try { $g->u .= []; } catch (Exception $e) { echo "|u:", gettype($g->u), $g->u; }
try { $g->m .= []; } catch (Exception $e) { echo "|m:", gettype($g->m), $g->m; }
restore_error_handler();
try { $g->f .= new stdClass; } catch (Error $e) { echo "|fo:", gettype($g->f), $g->f; }
