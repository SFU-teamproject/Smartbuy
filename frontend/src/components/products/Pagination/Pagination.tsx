import React from 'react';

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (pageNumber: number) => void;
    onPrevPage: () => void;
    onNextPage: () => void;
}

const Pagination: React.FC<PaginationProps> = ({
    currentPage,
    totalPages,
    onPageChange,
    onPrevPage,
    onNextPage,
}) => {
    const pageNumbers = Array.from({ length: totalPages }, (_, index) => index + 1);

    return (
        <div className='pagination-wrapper'>
            <div className='pagination-container'>
                <button onClick={onPrevPage} disabled={currentPage === 1} className='pagination-item'>
                    &lt;
                </button>

                {pageNumbers.map((number) => (
                    <button key={number} onClick={() => onPageChange(number)} className={`${currentPage === number ? 'pagination-item-active' : 'pagination-item'}`}>
                        {number}
                    </button>
                ))}

                <button onClick={onNextPage} disabled={currentPage === totalPages} className='pagination-item'>
                    &gt;
                </button>
            </div>
        </div>
    );
};

export default Pagination;