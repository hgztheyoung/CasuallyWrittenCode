(load "mk.scm")
(import (Framework match))
(import (Framework helpers))
(case-sensitive #t)
(print-gensym 'pretty)

(define-syntax λ
  (syntax-rules ()
    [(_ (x ...) body ...)
     (lambda (x ...) body ...)]))




#;
(define (cps1 program C)
  (match program
    [,sym (guard (symbol? sym)) (C sym)]
    [(λ (,x) ,body) (C `(λ (,x C) ,(cps1 body (λ (bcode) `(C ,bcode)))))]
    [(,app ,rator)
     (cps1 app (λ (acode)
                 (cps1 rator (λ (rcode)
                               (let ([sym (gensym)])
                                 `(,acode ,rcode (λ (,sym)
                                                   ,(C sym))))))))]))


(define (cps1 program C)
  (begin
    (display C)
    (display "\n")
   (match program
     [,sym (guard (symbol? sym)) (apply-cps C sym)]
     [(λ (,x) ,body) (apply-cps C `(λ (,x C) ,(cps1 body 'CONT-C0#;(λ (bcode) `(C ,bcode))
                                                    )))]
     [(,app ,rator)
      (cps1 app `(CONT-app-2 ,rator ,C) #;
            (λ (acode)
              (cps1 rator
                    `(CONT-app-final ,acode ,C)
                    #;
                    (λ (rcode)
                      (let ([sym (gensym)])
                        `(,acode ,rcode (λ (,sym)
                                          ,(apply-cps C sym))))))))])))


(define (apply-cps C-struct code)
  (match C-struct
    [CONT-ID code]
    [CONT-C0 `(C ,code)]
    [(CONT-app-final ,acode ,C)
     (let ([sym (gensym)])
       `(,acode ,code (λ (,sym)
                        ,(apply-cps C sym))))]
    [(CONT-app-2 ,rator ,C)
     (cps1 rator `(CONT-app-final ,code ,C))]))


#!eof

(load "defunc_cpser.ss")

(cps1 '((p q) (z c)) 'CONT-ID)
