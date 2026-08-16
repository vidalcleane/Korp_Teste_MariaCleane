import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatTableModule } from '@angular/material/table';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ProdutoService, Produto } from '../services/produto';

@Component({
  selector: 'app-produtos',
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatTableModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatCardModule,
    MatSnackBarModule
  ],
  templateUrl: './produtos.html',
  styleUrl: './produtos.css',
})
export class Produtos implements OnInit {
  colunas = ['id', 'codigo', 'descricao', 'saldo'];
  form: FormGroup;

  constructor(
    private produtoService: ProdutoService,
    private fb: FormBuilder,
    private snackBar: MatSnackBar
  ) {
    this.form = this.fb.group({
      codigo: ['', Validators.required],
      descricao: ['', Validators.required],
      saldo: [0, [Validators.required, Validators.min(0)]],
    });
  }

  // ngOnInit é um dos ciclos de vida do Angular: roda uma vez, assim que
  // o componente é montado na tela — o momento certo para buscar dados.
  ngOnInit(): void {
    this.produtoService.carregarProdutos();
  }

  get produtos$() {
    return this.produtoService.produtos$;
  }

  salvar(): void {
    if (this.form.invalid) {
      this.snackBar.open('Preencha todos os campos corretamente.', 'Fechar', { duration: 3000 });
      return;
    }

    const novoProduto: Produto = this.form.value;

    this.produtoService.criarProduto(novoProduto).subscribe({
      next: () => {
        this.snackBar.open('Produto cadastrado com sucesso!', 'Fechar', { duration: 3000 });
        this.form.reset({ codigo: '', descricao: '', saldo: 0 });
      },
      error: (err) => {
        const mensagem = err.error?.erro || 'Erro ao cadastrar produto.';
        this.snackBar.open(mensagem, 'Fechar', { duration: 4000 });
      }
    });
  }
}
