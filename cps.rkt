#lang racket
;seemingly right cps convert,with some referring to yinwang's cps.ss
#;(x -> (k x))
#;(((m n) -> (m n k)))
#;((m n) -> (m n (lambda (mnv) (k mnv))))
#;((x (y z)) -> (y z (lambda (yzres) (x yzres k))))
#;(((lambda (x) (m x)) (lambda (a) b)) -> ((lambda (x k0) (m x k0)) (lambda (a k1) (k1 b)) k))

(define (atom? expr)
  (match expr
    [x #:when (symbol? x) #t]
    [`(lambda (,x k) ,body) #t]
    [default #f]))


(define (cps expr)  
  (letrec ([id-gen (lambda (expr) `(k ,expr))]
           [cps1 (lambda (expr codegen)
                   (match expr
                     [atom #:when (atom? atom) (codegen expr)]
                     [`(,app ,rator) #:when (and (atom? app) (atom? rator))
                                     (if (eq? codegen id-gen)
                                         `(,app ,rator k) 
                                         (let ([x (gensym "var")])
                                           `(,app ,rator (lambda (,x) ,(codegen x)))))]
                     [`(,app ,rator) (cps1 rator (lambda (r-atom)
                                                   (cps1 app (lambda (a-atom)
                                                               (cps1 `(,a-atom ,r-atom) codegen)))))]
                     [`(lambda (,x) ,body) (codegen `(lambda (,x k) ,(cps1 body id-gen)))]))])
    (cps1 expr id-gen)))

(cps 'x)
(cps '(f x))
(cps '(x (y z)))
(cps '(lambda (x) (m x)))
(cps '((lambda (x) (m x)) (lambda (q) (lambda (a) b))))
'((lambda (x k) (m x k)) (lambda (q k) (k (lambda (a k) (k b)))) k)


(define (cps-without-eta-optimize expr)  
  (letrec ([id-gen (lambda (expr) `(k ,expr))]
           [cps1 (lambda (expr codegen)
                   (match expr
                     [atom #:when (atom? atom) (codegen expr)]
                     [`(,app ,rator) #:when (and (atom? app) (atom? rator))
                                     (let ([x (gensym "var")])
                                       `(,app ,rator (lambda (,x) ,(codegen x))))]
                     [`(,app ,rator) (cps1 rator (lambda (rpiece)
                                                   (cps1 app (lambda (apiece)
                                                               (cps1 `(,apiece ,rpiece) codegen)))))]
                     [`(lambda (,x) ,body) (codegen `(lambda (,x k) ,(cps1 body id-gen)))]))])
    (cps1 expr id-gen)))

(cps-without-eta-optimize '((lambda (x) (m x)) (lambda (q) (lambda (a) b))))
