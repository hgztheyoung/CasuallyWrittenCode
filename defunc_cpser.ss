(load "mk.scm")
(import (Framework match))
(import (Framework helpers))
(case-sensitive #t)
(print-gensym 'pretty)

(define-syntax λ
  (syntax-rules ()
    [(_ (x ...) body ...)
     (lambda (x ...) body ...)]))



;the original cps we all love
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


;defunctionalize it

(define (cps1 program C)
  (match program
    [,sym (guard (symbol? sym)) (apply-cps C sym)]
    [(λ (,x) ,body) (apply-cps C `(λ (,x C) ,(cps1 body 'CONT-C0)))]
    [(,app ,rator)
     (cps1 app `(CONT-app-2 ,rator ,C))]))

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


;we can make a defunctionalized program a relation
;thanks https://cgi.soic.indiana.edu/~c311/lib/exe/fetch.php?media=mk-convert.pdf

(define (caro p a)
  (fresh (d)
    (== p (cons a d))))

(define (cps1ᵒ program C out)
  (conde
   [(symbolo program) (apply-cpsᵒ C program out)]
   [(fresh (x body bodycode)
      (== program `(λ (,x) ,body))
      (symbolo x)      
      (apply-cpsᵒ C `(λ (,x C) ,bodycode) out)
      (cps1ᵒ body 'CONT-C0 bodycode))]
   [(fresh (app rator)
      (== program `(,app ,rator))
      (cps1ᵒ app `(CONT-app-2 ,rator ,C) out))]))


(define (apply-cpsᵒ C-struct code out)
  (conde
   [(== C-struct 'CONT-ID) (== code out)]
   [(== C-struct 'CONT-C0) (== `(C ,code) out)]
   [(caro C-struct 'CONT-app-final)
    (fresh (sym appcode acode C)
      (== C-struct `(CONT-app-final ,acode ,C))
      (== `(,acode ,code (λ (,sym) ,appcode)) out)
      (symbolo sym)
      (apply-cpsᵒ C sym appcode))]
   [(caro C-struct 'CONT-app-2)
    (fresh (rator C)
      (== C-struct `(CONT-app-2 ,rator ,C))
      (cps1ᵒ rator `(CONT-app-final ,code ,C) out))]))




#!eof

(load "script.ss")

;(cps1 '((p q) (z c)) 'CONT-ID)

(run 1 (q)
  (cps1ᵒ '(z ((p q) (z (s c)))) 'CONT-ID q))

(run 1 (q)
  (cps1ᵒ q 'CONT-ID
   '(p q (λ (_.0) (z c (λ (_.1) (_.0 _.1 (λ (_.2) _.2))))))))

(run 1 (q)
  (cps1ᵒ 'z q
   '(p q
       (λ (_.0)
         (s c
            (λ (_.1)
              (z _.1
                 (λ (_.2) (_.0 _.2 (λ (_.3) (z _.3 (λ (_.4) _.4))))))))))))



(run 1 (q)
  (cps1ᵒ '((p q) (z c)) q
   '(p q (λ (_.0) (z c (λ (_.1) (_.0 _.1 (λ (_.2) _.2))))))))


(run 1 (q)
  (cps1ᵒ '(p q) q
   '(p q (λ (_.0) (z c (λ (_.1) (_.0 _.1 (λ (_.2) _.2))))))))


;fail to terminate
#;
(run 2 (q)
  (cps1ᵒ '(p q) q
         '(p q (λ (_.0) (z c (λ (_.1) (_.0 _.1 (λ (_.2) _.2))))))))



(run 1000 (p k q)
  (cps1ᵒ p k q))
